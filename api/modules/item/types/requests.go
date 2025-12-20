package types

type CreateItemRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateItemRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Completed   bool    `json:"completed"`
}
