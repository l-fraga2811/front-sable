package supabase

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// jwksResponse representa a resposta do endpoint JWKS.
//
// CONCEITO: JWKS (JSON Web Key Set)
// JWKS é um padrão para publicar chaves públicas usadas para verificar
// assinaturas de tokens JWT. O endpoint JWKS retorna um JSON com um
// array de chaves públicas.
//
// URL típica do Supabase: https://xxx.supabase.co/auth/v1/.well-known/jwks.json
type jwksResponse struct {
	Keys []jwkKey `json:"keys"`
}

// jwkKey representa uma chave pública no formato JWK.
//
// CAMPOS IMPORTANTES:
//   - kty: Tipo da chave (RSA, EC, etc)
//   - kid: ID único da chave (usado para encontrar a chave certa)
//   - n: Módulo RSA (parte da chave pública)
//   - e: Expoente RSA (parte da chave pública)
//   - alg: Algoritmo (RS256, RS384, etc)
//   - use: Uso da chave (sig = assinatura, enc = encriptação)
type jwkKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JwksCache é um cache thread-safe para chaves públicas JWKS.
//
// POR QUE CACHE?
// Buscar as chaves do endpoint JWKS a cada requisição seria lento e
// poderia causar rate limiting. O cache armazena as chaves por um
// período (TTL) e só busca novamente quando expiram.
//
// THREAD-SAFETY:
// Usamos sync.RWMutex para permitir múltiplas leituras simultâneas
// (RLock) mas apenas uma escrita por vez (Lock).
type JwksCache struct {
	jwksURL string
	client  *http.Client

	mu        sync.RWMutex
	publicKey map[string]*rsa.PublicKey
	expiresAt time.Time
	ttl       time.Duration
}

// NewJwksCache cria um novo cache de chaves JWKS.
func NewJwksCache(jwksURL string) *JwksCache {
	return &JwksCache{
		jwksURL:   jwksURL,
		client:    &http.Client{Timeout: 10 * time.Second},
		publicKey: map[string]*rsa.PublicKey{},
		ttl:       10 * time.Minute,
	}
}

// GetPublicKey retorna a chave pública para um dado kid (key ID).
//
// FLUXO:
//  1. Tenta ler do cache (RLock - permite leituras concorrentes)
//  2. Se não encontrar ou cache expirado, faz refresh
//  3. Retorna a chave ou erro se não encontrada
func (c *JwksCache) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	if kid == "" {
		return nil, errors.New("kid não informado")
	}

	c.mu.RLock()
	if time.Now().Before(c.expiresAt) {
		if key, ok := c.publicKey[kid]; ok {
			c.mu.RUnlock()
			return key, nil
		}
	}
	c.mu.RUnlock()

	if err := c.refresh(); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	key, ok := c.publicKey[kid]
	if !ok {
		return nil, errors.New("kid não encontrado no jwks")
	}
	return key, nil
}

// refresh busca as chaves do endpoint JWKS e atualiza o cache.
//
// PADRÃO: Double-Check Locking
// Verificamos se o cache expirou DUAS vezes:
//  1. Antes de adquirir o lock (evita lock desnecessário)
//  2. Depois de adquirir o lock (evita refresh duplicado)
func (c *JwksCache) refresh() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().Before(c.expiresAt) {
		return nil
	}

	resp, err := c.client.Get(c.jwksURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("falha ao buscar jwks")
	}

	var parsed jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return err
	}

	newKeys := make(map[string]*rsa.PublicKey, len(parsed.Keys))
	for _, key := range parsed.Keys {
		if key.Kty != "RSA" || key.N == "" || key.E == "" || key.Kid == "" {
			continue
		}

		pub, err := rsaFromJwk(key.N, key.E)
		if err != nil {
			continue
		}
		newKeys[key.Kid] = pub
	}

	if len(newKeys) == 0 {
		return errors.New("jwks sem chaves válidas")
	}

	c.publicKey = newKeys
	c.expiresAt = time.Now().Add(c.ttl)
	return nil
}

// rsaFromJwk converte os componentes JWK (n, e) para uma chave RSA pública.
//
// MATEMÁTICA RSA:
//   - n (módulo): produto de dois primos grandes
//   - e (expoente): geralmente 65537 (0x10001)
//
// Ambos são codificados em Base64URL no JWK.
func rsaFromJwk(nB64 string, eB64 string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)
	if !e.IsInt64() {
		return nil, errors.New("expoente inválido")
	}

	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}
