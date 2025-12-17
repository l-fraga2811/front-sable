// Package handlers contém os handlers HTTP da aplicação.
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-api/middleware"
	"go-api/models"
	"go-api/storage"
)

// ItemHandler agrupa os handlers do CRUD de items.
type ItemHandler struct {
	store *storage.Storage
}

// NewItemHandler cria uma nova instância de ItemHandler.
func NewItemHandler() *ItemHandler {
	return &ItemHandler{
		store: storage.GetStorage(),
	}
}

// Create cria um novo item.
// @Summary Criar item
// @Description Cria um novo item para o usuário autenticado
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateItemRequest true "Dados do item"
// @Success 201 {object} models.Item
// @Failure 400 {object} map[string]string
// @Router /api/items [post]
func (h *ItemHandler) Create(c *gin.Context) {
	// Obtém o ID do usuário do contexto (definido pelo middleware de auth).
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	// Faz o parse do body JSON para a struct de request.
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Cria o item com os dados fornecidos.
	// time.Now() retorna a data/hora atual.
	item := models.Item{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      userID,
	}

	// Salva o item no storage.
	if err := h.store.CreateItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao criar item",
		})
		return
	}

	// Retorna o item criado com status 201 (Created).
	c.JSON(http.StatusCreated, item)
}

// GetAll retorna todos os items do usuário autenticado.
// @Summary Listar items
// @Description Retorna todos os items do usuário autenticado
// @Tags items
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Item
// @Router /api/items [get]
func (h *ItemHandler) GetAll(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	// Busca todos os items do usuário.
	items := h.store.GetItemsByUserID(userID)

	// Se não houver items, retorna um array vazio (não null).
	// Isso é uma boa prática para APIs REST.
	if items == nil {
		items = []models.Item{}
	}

	c.JSON(http.StatusOK, items)
}

// GetByID retorna um item específico pelo ID.
// @Summary Buscar item por ID
// @Description Retorna um item específico do usuário autenticado
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID do item"
// @Success 200 {object} models.Item
// @Failure 404 {object} map[string]string
// @Router /api/items/{id} [get]
func (h *ItemHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	// c.Param obtém parâmetros da URL.
	// Para a rota "/api/items/:id", c.Param("id") retorna o valor de :id.
	itemID := c.Param("id")

	// Busca o item pelo ID.
	item, err := h.store.GetItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item não encontrado",
		})
		return
	}

	// Verifica se o item pertence ao usuário autenticado.
	// Isso é importante para segurança - um usuário não pode ver items de outros.
	if item.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Você não tem permissão para acessar este item",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Update atualiza um item existente.
// @Summary Atualizar item
// @Description Atualiza um item do usuário autenticado
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID do item"
// @Param request body models.UpdateItemRequest true "Dados para atualizar"
// @Success 200 {object} models.Item
// @Failure 404 {object} map[string]string
// @Router /api/items/{id} [put]
func (h *ItemHandler) Update(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	itemID := c.Param("id")

	// Busca o item existente.
	item, err := h.store.GetItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item não encontrado",
		})
		return
	}

	// Verifica permissão.
	if item.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Você não tem permissão para atualizar este item",
		})
		return
	}

	// Faz o parse dos dados de atualização.
	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Atualiza apenas os campos fornecidos.
	// Em Go, strings vazias e zeros são os "zero values".
	// Você pode usar ponteiros para diferenciar "não enviado" de "enviado vazio".
	if req.Title != "" {
		item.Title = req.Title
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Price != 0 {
		item.Price = req.Price
	}
	// Completed é um bool, então sempre atualizamos.
	item.Completed = req.Completed
	item.UpdatedAt = time.Now()

	// Salva as alterações.
	if err := h.store.UpdateItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar item",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Delete remove um item.
// @Summary Deletar item
// @Description Remove um item do usuário autenticado
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID do item"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/items/{id} [delete]
func (h *ItemHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	itemID := c.Param("id")

	// Busca o item para verificar permissão.
	item, err := h.store.GetItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item não encontrado",
		})
		return
	}

	// Verifica permissão.
	if item.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Você não tem permissão para deletar este item",
		})
		return
	}

	// Remove o item.
	if err := h.store.DeleteItem(itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao deletar item",
		})
		return
	}

	// Retorna mensagem de sucesso.
	// Alguns preferem retornar 204 (No Content) sem body.
	c.JSON(http.StatusOK, gin.H{
		"message": "Item deletado com sucesso",
	})
}
