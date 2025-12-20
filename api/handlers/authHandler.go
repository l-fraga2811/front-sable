// Package handlers contém os handlers HTTP da aplicação.
//
// CONCEITO: Handlers (ou Controllers)
// Handlers são funções que processam requisições HTTP e retornam respostas.
// No padrão MVC, eles seriam os "Controllers".
// Cada handler é responsável por:
//   - Receber e validar dados da requisição
//   - Chamar a lógica de negócio (services/repositories)
//   - Formatar e retornar a resposta
//
// PADRÃO: Struct com métodos
// Agrupamos handlers relacionados em uma struct para:
//   - Compartilhar dependências (ex: conexão com banco)
//   - Organizar código por domínio (auth, items, etc)
//   - Facilitar testes (injeção de dependências)
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api/middleware"
	"go-api/models"
)

// AuthHandler agrupa os handlers relacionados à autenticação.
//
// NOTA: Com Supabase Auth, a autenticação (login/register/logout) é feita
// diretamente no frontend usando o SDK do Supabase. Este handler agora
// serve apenas para retornar dados do usuário autenticado.
type AuthHandler struct{}

// NewAuthHandler cria uma nova instância de AuthHandler.
//
// PADRÃO GO: Constructor Function
// Go não tem construtores como Java/C#. Usamos funções que retornam
// ponteiros para a struct inicializada. Por convenção, o nome é New + NomeDaStruct.
//
// Retornamos ponteiro (*AuthHandler) para:
//   - Evitar cópia da struct em cada chamada
//   - Permitir que métodos modifiquem o estado (se necessário)
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// GetProfile retorna os dados do usuário autenticado.
//
// FLUXO:
//  1. O middleware de autenticação já validou o token JWT
//  2. O userID e email foram extraídos do token e salvos no contexto
//  3. Este handler apenas lê esses dados e retorna
//
// @Summary Obter perfil do usuário
// @Description Retorna os dados do usuário autenticado via token JWT
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]string "Usuário não autenticado"
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	email, _ := c.Get("email")
	emailStr, _ := email.(string)

	c.JSON(http.StatusOK, models.UserResponse{
		ID:       userID,
		Username: "",
		Email:    emailStr,
	})
}
