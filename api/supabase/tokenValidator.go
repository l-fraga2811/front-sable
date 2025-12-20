package supabase

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenValidator valida tokens JWT do Supabase Auth.
//
// CONCEITO: JWT (JSON Web Token)
// JWT é um padrão para transmitir informações de forma segura entre partes.
// Um JWT tem 3 partes separadas por ".":
//   - Header: algoritmo de assinatura (HS256, RS256, etc)
//   - Payload: dados (claims) como userID, email, expiração
//   - Signature: assinatura para verificar autenticidade
//
// ALGORITMOS SUPORTADOS:
//   - HS256: Assinatura simétrica (usa um segredo compartilhado)
//   - RS256: Assinatura assimétrica (usa par de chaves pública/privada)
//
// O Supabase pode usar ambos dependendo da configuração do projeto.
type TokenValidator struct {
	jwks      *JwksCache
	jwtSecret []byte
}

// NewTokenValidator cria um novo validador de tokens.
func NewTokenValidator(cfg Config) *TokenValidator {
	return &TokenValidator{
		jwks:      NewJwksCache(cfg.JwksURL),
		jwtSecret: []byte(cfg.JwtSecret),
	}
}

// AccessTokenClaims representa os dados extraídos do token JWT.
//
// CONCEITO: Claims
// Claims são as "afirmações" contidas no token. Existem claims padrão
// (RegisteredClaims) e claims customizados. O Supabase inclui:
//   - sub: ID do usuário (Subject)
//   - email: Email do usuário
//   - role: Papel do usuário (authenticated, anon, etc)
//   - exp: Data de expiração
//   - iat: Data de emissão
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Validate verifica se um token JWT é válido e retorna seus claims.
//
// FLUXO DE VALIDAÇÃO:
//  1. Extrai o header do token para descobrir o algoritmo (alg)
//  2. Se HS256: valida usando SUPABASE_JWT_SECRET
//  3. Se RS256: valida usando chave pública do JWKS
//  4. Verifica se o token não expirou
//  5. Verifica se tem o campo "sub" (userID)
//
// SEGURANÇA: Sempre validamos o algoritmo explicitamente para evitar
// ataques de "algorithm confusion" onde um atacante força uso de "none".
func (v *TokenValidator) Validate(tokenString string) (AccessTokenClaims, error) {
	alg, kid, err := tokenHeader(tokenString)
	if err != nil {
		return AccessTokenClaims{}, err
	}

	claims := AccessTokenClaims{}
	var parsed *jwt.Token

	if alg == "HS256" {
		if len(v.jwtSecret) == 0 {
			return AccessTokenClaims{}, errors.New("SUPABASE_JWT_SECRET não definido")
		}
		parsed, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return v.jwtSecret, nil
		}, jwt.WithValidMethods([]string{"HS256"}))
	} else {
		parsed, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			key, err := v.jwks.GetPublicKey(kid)
			if err != nil {
				return nil, err
			}
			return key, nil
		}, jwt.WithValidMethods([]string{"RS256"}))
	}

	if err != nil {
		return AccessTokenClaims{}, err
	}
	if parsed == nil || !parsed.Valid {
		return AccessTokenClaims{}, errors.New("token inválido")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return AccessTokenClaims{}, errors.New("token expirado")
	}

	if claims.Subject == "" {
		return AccessTokenClaims{}, errors.New("sub ausente")
	}

	return claims, nil
}

// tokenHeader extrai o algoritmo (alg) e key ID (kid) do header do JWT.
//
// ESTRUTURA DO JWT:
// header.payload.signature
//
// O header é codificado em Base64URL e contém:
//
//	{"alg": "HS256", "typ": "JWT", "kid": "..."}
func tokenHeader(tokenString string) (string, string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", "", errors.New("token mal formatado")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", "", err
	}

	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return "", "", err
	}
	if header.Alg == "" {
		return "", "", errors.New("alg ausente")
	}
	if header.Alg != "HS256" && header.Kid == "" {
		return "", "", errors.New("kid ausente")
	}

	return header.Alg, header.Kid, nil
}
