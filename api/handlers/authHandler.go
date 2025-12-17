// Package handlers contém os handlers HTTP da aplicação.
// Handlers são funções que processam requests HTTP e retornam responses.
// Cada handler é responsável por uma operação específica da API.
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go-api/middleware"
	"go-api/models"
	"go-api/storage"
)

// AuthHandler agrupa os handlers de autenticação.
// Usar uma struct permite compartilhar dependências entre handlers.
type AuthHandler struct {
	// store é a referência ao storage.
	// Em uma aplicação real, isso seria uma interface de repositório.
	store *storage.Storage
}

// NewAuthHandler cria uma nova instância de AuthHandler.
// Este padrão é chamado de "constructor function" em Go.
// Go não tem construtores como outras linguagens, então usamos funções.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		store: storage.GetStorage(),
	}
}

// Register cria um novo usuário.
// @Summary Registrar novo usuário
// @Description Cria uma nova conta de usuário
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Dados do usuário"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// Cria uma variável para armazenar os dados do request.
	var req models.RegisterRequest

	// ShouldBindJSON faz o parse do JSON do body para a struct.
	// Se houver erro (JSON inválido ou campos obrigatórios faltando),
	// retorna um erro.
	if err := c.ShouldBindJSON(&req); err != nil {
		// gin.H é um atalho para map[string]interface{}.
		// É muito usado para criar responses JSON rapidamente.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Verifica se o usuário já existe.
	_, err := h.store.GetUserByUsername(req.Username)
	if err == nil {
		// Se não houve erro, significa que o usuário existe.
		c.JSON(http.StatusConflict, gin.H{
			"error": "Usuário já existe",
		})
		return
	}

	// Gera o hash da senha usando bcrypt.
	// bcrypt é um algoritmo de hash seguro para senhas.
	// O segundo parâmetro é o "cost" - quanto maior, mais seguro (e mais lento).
	// bcrypt.DefaultCost é 10, que é um bom equilíbrio.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao processar senha",
		})
		return
	}

	// Cria o usuário com um ID único.
	// uuid.New() gera um UUID v4 (aleatório).
	user := models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// Salva o usuário no storage.
	if err := h.store.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao criar usuário",
		})
		return
	}

	// Retorna o usuário criado (sem a senha).
	// http.StatusCreated (201) indica que um recurso foi criado.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuário criado com sucesso",
		"user":    user.ToResponse(),
	})
}

// Login autentica um usuário e retorna um token JWT.
// @Summary Login de usuário
// @Description Autentica o usuário e retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Credenciais"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Busca o usuário pelo username.
	user, err := h.store.GetUserByUsername(req.Username)
	if err != nil {
		// Retornamos uma mensagem genérica por segurança.
		// Não revelamos se o usuário existe ou não.
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Credenciais inválidas",
		})
		return
	}

	// Compara a senha fornecida com o hash armazenado.
	// bcrypt.CompareHashAndPassword retorna nil se as senhas coincidem.
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Credenciais inválidas",
		})
		return
	}

	// Cria os claims do token JWT.
	// O token expira em 24 horas.
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// ExpiresAt define quando o token expira.
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// IssuedAt define quando o token foi criado.
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Subject é geralmente o ID do usuário.
			Subject: user.ID,
		},
	}

	// Cria o token JWT.
	// jwt.SigningMethodHS256 é o algoritmo de assinatura (HMAC-SHA256).
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a chave secreta.
	// SignedString retorna o token como string.
	tokenString, err := token.SignedString(middleware.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao gerar token",
		})
		return
	}

	// Retorna o token e informações do usuário.
	c.JSON(http.StatusOK, gin.H{
		"message":   "Login realizado com sucesso",
		"token":     tokenString,
		"expiresAt": expirationTime.Format(time.RFC3339),
		"user":      user.ToResponse(),
	})
}

// Logout invalida o token atual.
// @Summary Logout de usuário
// @Description Invalida o token JWT atual
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Obtém o token do contexto (foi salvo pelo middleware).
	token := middleware.GetTokenFromContext(c)

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token não encontrado",
		})
		return
	}

	// Adiciona o token à blacklist.
	// Isso impede que o token seja usado novamente.
	h.store.BlacklistToken(token)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout realizado com sucesso",
	})
}

// GetProfile retorna os dados do usuário autenticado.
// @Summary Obter perfil do usuário
// @Description Retorna os dados do usuário autenticado
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Obtém o username do contexto (foi salvo pelo middleware).
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	// Busca o usuário no storage.
	user, err := h.store.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuário não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
