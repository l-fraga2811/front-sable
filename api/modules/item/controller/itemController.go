package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-api/middleware"
	"go-api/models"
	"go-api/modules/item/types"
	"go-api/supabase"
)

// ItemController agrupa a lógica de CRUD de items.
type ItemController struct {
	client *supabase.RestClient
}

// NewItemController cria o controller com a dependência do Supabase.
func NewItemController() *ItemController {
	client, err := supabase.GetClient()
	if err != nil {
		return &ItemController{}
	}
	return &ItemController{
		client: client.RestClient,
	}
}

func (ctrl *ItemController) Create(c *gin.Context) {
	userID, accessToken, ok := ctrl.requireAuth(c)
	if !ok {
		return
	}

	var req types.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	created, err := ctrl.client.CreateItem(c.Request.Context(), accessToken, supabase.CreateItemPayload{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar item"})
		return
	}

	item, err := mapItemRow(created)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (ctrl *ItemController) GetAll(c *gin.Context) {
	_, accessToken, ok := ctrl.requireAuth(c)
	if !ok {
		return
	}
	if ctrl.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cliente Supabase não inicializado"})
		return
	}

	rows, err := ctrl.client.ListItems(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar items"})
		return
	}

	items := make([]models.Item, 0, len(rows))
	for _, row := range rows {
		item, err := mapItemRow(row)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar items"})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

func (ctrl *ItemController) GetByID(c *gin.Context) {
	userID, accessToken, ok := ctrl.requireAuth(c)
	if !ok {
		return
	}
	if ctrl.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cliente Supabase não inicializado"})
		return
	}

	itemID := c.Param("id")
	row, found, err := ctrl.client.GetItemByID(c.Request.Context(), accessToken, itemID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar item"})
		return
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	item, err := mapItemRow(row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar item"})
		return
	}
	if item.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para acessar este item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (ctrl *ItemController) Update(c *gin.Context) {
	userID, accessToken, ok := ctrl.requireAuth(c)
	if !ok {
		return
	}
	if ctrl.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cliente Supabase não inicializado"})
		return
	}

	itemID := c.Param("id")
	row, found, err := ctrl.client.GetItemByID(c.Request.Context(), accessToken, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar item"})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	item, err := mapItemRow(row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar item"})
		return
	}
	if item.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para atualizar este item"})
		return
	}

	var req types.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	payload := supabase.UpdateItemPayload{}
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

	updatedRow, updated, err := ctrl.client.UpdateItem(c.Request.Context(), accessToken, itemID, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar item"})
		return
	}
	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	updatedItem, err := mapItemRow(updatedRow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar item"})
		return
	}
	if updatedItem.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para atualizar este item"})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

func (ctrl *ItemController) Delete(c *gin.Context) {
	_, accessToken, ok := ctrl.requireAuth(c)
	if !ok {
		return
	}
	if ctrl.client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cliente Supabase não inicializado"})
		return
	}

	itemID := c.Param("id")
	deleted, err := ctrl.client.DeleteItem(c.Request.Context(), accessToken, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar item"})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deletado com sucesso"})
}

func (ctrl *ItemController) requireAuth(c *gin.Context) (string, string, bool) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return "", "", false
	}

	accessToken := middleware.GetTokenFromContext(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return "", "", false
	}

	return userID, accessToken, true
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
