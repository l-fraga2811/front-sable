// Package handlers contém os handlers HTTP da aplicação.
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-api/middleware"
	"go-api/models"
	"go-api/supabase"
)

// ItemHandler agrupa os handlers do CRUD de items.
type ItemHandler struct {
	client *supabase.RestClient
}

// NewItemHandler cria uma nova instância de ItemHandler.
func NewItemHandler() *ItemHandler {
	client, err := supabase.GetClient()
	if err != nil {
		return &ItemHandler{}
	}
	return &ItemHandler{
		client: client.RestClient,
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

	accessToken := middleware.GetTokenFromContext(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}
	if h.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cliente Supabase não inicializado",
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

	created, err := h.client.CreateItem(c.Request.Context(), accessToken, supabase.CreateItemPayload{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Completed:   false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao criar item",
		})
		return
	}

	item, err := mapItemRow(created)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao processar item",
		})
		return
	}

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

	accessToken := middleware.GetTokenFromContext(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}
	if h.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cliente Supabase não inicializado",
		})
		return
	}

	rows, err := h.client.ListItems(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar items",
		})
		return
	}

	items := make([]models.Item, 0, len(rows))
	for _, row := range rows {
		item, err := mapItemRow(row)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao processar items",
			})
			return
		}
		items = append(items, item)
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

	accessToken := middleware.GetTokenFromContext(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}
	if h.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cliente Supabase não inicializado",
		})
		return
	}

	// c.Param obtém parâmetros da URL.
	// Para a rota "/api/items/:id", c.Param("id") retorna o valor de :id.
	itemID := c.Param("id")

	row, found, err := h.client.GetItemByID(c.Request.Context(), accessToken, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar item",
		})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item não encontrado",
		})
		return
	}

	item, err := mapItemRow(row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao processar item",
		})
		return
	}
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

	accessToken := middleware.GetTokenFromContext(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}
	if h.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cliente Supabase não inicializado",
		})
		return
	}

	itemID := c.Param("id")

	row, found, err := h.client.GetItemByID(c.Request.Context(), accessToken, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar item",
		})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item não encontrado",
		})
		return
	}

	item, err := mapItemRow(row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao processar item",
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

	var payload supabase.UpdateItemPayload
	if req.Title != "" {
		payload.Title = &req.Title
	}
	if req.Description != "" {
		payload.Description = &req.Description
	}
	if req.Price != 0 {
		payload.Price = &req.Price
	}
	payload.Completed = &req.Completed

	updatedRow, updated, err := h.client.UpdateItem(c.Request.Context(), accessToken, itemID, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar item",
		})
		return
	}
	if !updated {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item não encontrado",
		})
		return
	}

	updatedItem, err := mapItemRow(updatedRow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao processar item",
		})
		return
	}

	if updatedItem.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Você não tem permissão para atualizar este item",
		})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
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

	accessToken := middleware.GetTokenFromContext(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}
	if h.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cliente Supabase não inicializado",
		})
		return
	}

	itemID := c.Param("id")

	deleted, err := h.client.DeleteItem(c.Request.Context(), accessToken, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao deletar item",
		})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Item não encontrado",
		})
		return
	}

	// Retorna mensagem de sucesso.
	// Alguns preferem retornar 204 (No Content) sem body.
	c.JSON(http.StatusOK, gin.H{
		"message": "Item deletado com sucesso",
	})
}

func mapItemRow(row supabase.ItemRow) (models.Item, error) {
	createdAt, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
	if err != nil {
		return models.Item{}, err
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	if err != nil {
		return models.Item{}, err
	}

	return models.Item{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Price:       float64(row.Price),
		Completed:   row.Completed,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		UserID:      row.UserID,
	}, nil
}
