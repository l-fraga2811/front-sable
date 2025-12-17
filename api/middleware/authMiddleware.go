// Package middleware contém os middlewares da aplicação.
// Middlewares são funções que executam ANTES dos handlers.
// São úteis para: autenticação, logging, CORS, rate limiting, etc.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go-api/storage"
)

// JWTSecret é a chave secreta usada para assinar os tokens JWT.
// IMPORTANTE: Em produção, NUNCA deixe isso hardcoded!
// Use variáveis de ambiente: os.Getenv("JWT_SECRET")
var JWTSecret = []byte("minha-chave-secreta-super-segura-mude-em-producao")

// Claims é a estrutura que define o "payload" do token JWT.
// Ela "embute" jwt.RegisteredClaims para ter os campos padrão do JWT.
type Claims struct {
	// UserID é o ID do usuário autenticado.
	UserID string `json:"userId"`

	// Username é o nome do usuário.
	Username string `json:"username"`

	// jwt.RegisteredClaims contém campos padrão como:
	// - ExpiresAt: quando o token expira
	// - IssuedAt: quando o token foi criado
	// - Subject: assunto do token
	jwt.RegisteredClaims
}

// AuthMiddleware é o middleware que protege rotas que requerem autenticação.
// Ele verifica se o token JWT é válido antes de permitir acesso à rota.
func AuthMiddleware() gin.HandlerFunc {
	// Retornamos uma função que será executada para cada request.
	// gin.HandlerFunc é o tipo de função que o Gin espera para handlers/middlewares.
	return func(c *gin.Context) {
		// Obtém o header "Authorization" do request.
		// O formato esperado é: "Bearer <token>"
		authHeader := c.GetHeader("Authorization")

		// Verifica se o header existe e começa com "Bearer ".
		if authHeader == "" {
			// c.JSON envia uma resposta JSON.
			// c.Abort() impede que os próximos handlers sejam executados.
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token de autorização não fornecido",
			})
			c.Abort()
			return
		}

		// Verifica se o header tem o prefixo "Bearer ".
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Extrai o token removendo o prefixo "Bearer ".
		// strings.TrimPrefix remove o prefixo se existir.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verifica se o token está na blacklist (foi feito logout).
		store := storage.GetStorage()
		if store.IsTokenBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token foi invalidado (logout realizado)",
			})
			c.Abort()
			return
		}

		// Faz o parse e validação do token JWT.
		// jwt.ParseWithClaims decodifica o token e valida a assinatura.
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Esta função retorna a chave secreta para validar a assinatura.
			// Você pode adicionar validações extras aqui, como verificar o algoritmo.
			return JWTSecret, nil
		})

		// Verifica se houve erro no parse ou se o token é inválido.
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido ou expirado",
			})
			c.Abort()
			return
		}

		// Se chegou aqui, o token é válido!
		// Salvamos os dados do usuário no contexto para uso nos handlers.
		// c.Set salva um valor que pode ser recuperado com c.Get ou c.MustGet.
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("token", tokenString)

		// c.Next() passa o controle para o próximo handler na cadeia.
		// Se não chamarmos Next(), os handlers seguintes não executam.
		c.Next()
	}
}

// GetUserIDFromContext é uma função helper para extrair o userID do contexto.
// Isso evita repetição de código nos handlers.
func GetUserIDFromContext(c *gin.Context) string {
	// MustGet retorna o valor ou causa panic se não existir.
	// Use Get se quiser tratar o caso de não existir.
	userID, exists := c.Get("userID")
	if !exists {
		return ""
	}
	// Type assertion: converte interface{} para string.
	// O segundo valor (ok) indica se a conversão foi bem sucedida.
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}

// GetTokenFromContext extrai o token do contexto.
func GetTokenFromContext(c *gin.Context) string {
	token, exists := c.Get("token")
	if !exists {
		return ""
	}
	if t, ok := token.(string); ok {
		return t
	}
	return ""
}
