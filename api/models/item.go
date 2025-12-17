// Package models contém as estruturas de dados da aplicação.
package models

import "time"

// Item representa um item genérico para o CRUD.
// Esta struct demonstra vários tipos de dados comuns em Go.
type Item struct {
	// ID é o identificador único do item.
	// Usamos string para simplificar, mas em produção você usaria UUID.
	ID string `json:"id"`

	// Title é o título do item.
	Title string `json:"title" binding:"required"`

	// Description é a descrição do item.
	// O ponteiro (*string) permite que o campo seja nil (nulo).
	// Isso é útil para campos opcionais.
	Description string `json:"description"`

	// Price demonstra o uso de float64 para valores decimais.
	// Em produção, para dinheiro, use bibliotecas específicas para evitar
	// problemas de precisão com ponto flutuante.
	Price float64 `json:"price"`

	// Completed é um booleano que indica se o item está completo.
	Completed bool `json:"completed"`

	// CreatedAt armazena quando o item foi criado.
	// O pacote "time" é a biblioteca padrão de Go para datas.
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt armazena quando o item foi atualizado pela última vez.
	UpdatedAt time.Time `json:"updatedAt"`

	// UserID é o ID do usuário que criou o item.
	// Isso permite que cada usuário tenha seus próprios items.
	UserID string `json:"userId"`
}

// CreateItemRequest é a struct para criar um novo item.
// Separamos request de response para ter controle sobre quais campos
// o cliente pode enviar vs quais campos retornamos.
type CreateItemRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// UpdateItemRequest é a struct para atualizar um item.
// Todos os campos são opcionais (sem binding:"required").
type UpdateItemRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Completed   bool    `json:"completed"`
}
