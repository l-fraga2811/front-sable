// Package middleware contém os middlewares da aplicação.
// Middlewares são funções que executam ANTES dos handlers.
// São úteis para: autenticação, logging, CORS, rate limiting, etc.
package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-api/supabase"
)

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

		client, err := supabase.GetClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao inicializar Supabase",
			})
			c.Abort()
			return
		}

		claims, err := client.TokenValidator.Validate(tokenString)
		if err != nil {
			log.Printf("[AUTH] Erro ao validar token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido ou expirado",
			})
			c.Abort()
			return
		}

		// Se chegou aqui, o token é válido!
		// Salvamos os dados do usuário no contexto para uso nos handlers.
		// c.Set salva um valor que pode ser recuperado com c.Get ou c.MustGet.
		c.Set("userID", claims.Subject)
		c.Set("email", claims.Email)
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
