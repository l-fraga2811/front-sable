// Package supabase contém a integração com o Supabase (Auth + Database).
//
// ARQUITETURA DO PACOTE:
//   - config.go     → Carrega configurações do ambiente
//   - client.go     → Cliente singleton que agrupa todos os componentes
//   - tokenValidator.go → Valida tokens JWT (HS256 e RS256)
//   - jwks.go       → Cache de chaves públicas para validação RS256
//   - restClient.go → Cliente HTTP para a API REST do Supabase (PostgREST)
//   - numeric.go    → Tipo customizado para campos numeric do Postgres
package supabase

import (
	"errors"
	"os"
	"strings"
)

// Config armazena as configurações necessárias para conectar ao Supabase.
//
// VARIÁVEIS DE AMBIENTE:
//   - SUPABASE_PROJECT_URL: URL do projeto (ex: https://xxx.supabase.co)
//   - SUPABASE_ANON_KEY: Chave pública (anon key) para autenticação
//   - SUPABASE_JWT_SECRET: Segredo para validar tokens HS256
//   - SUPABASE_JWKS_URL: (opcional) URL do JWKS para validar tokens RS256
type Config struct {
	ProjectURL string
	AnonKey    string
	JwksURL    string
	JwtSecret  string
}

// LoadConfig carrega as configurações do Supabase a partir das variáveis de ambiente.
//
// PADRÃO: Fail Fast
// Se uma configuração obrigatória estiver faltando, retornamos erro imediatamente.
// Isso evita erros confusos em tempo de execução.
//
// PADRÃO: Fallback
// Para SUPABASE_JWKS_URL, derivamos automaticamente da URL do projeto se não for
// explicitamente definida. Isso reduz a configuração necessária.
func LoadConfig() (Config, error) {
	projectURL := strings.TrimRight(os.Getenv("SUPABASE_PROJECT_URL"), "/")
	if projectURL == "" {
		return Config{}, errors.New("SUPABASE_PROJECT_URL não definido")
	}

	anonKey := os.Getenv("SUPABASE_ANON_KEY")
	if anonKey == "" {
		anonKey = os.Getenv("SUPABASE_API_KEY")
	}
	if anonKey == "" {
		return Config{}, errors.New("SUPABASE_ANON_KEY (ou SUPABASE_API_KEY) não definido")
	}

	jwksURL := os.Getenv("SUPABASE_JWKS_URL")
	if jwksURL == "" {
		jwksURL = projectURL + "/auth/v1/.well-known/jwks.json"
	}

	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")

	return Config{
		ProjectURL: projectURL,
		AnonKey:    anonKey,
		JwksURL:    jwksURL,
		JwtSecret:  jwtSecret,
	}, nil
}
