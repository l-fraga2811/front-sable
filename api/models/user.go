// Package models contém as estruturas de dados (structs) da aplicação.
// Em Go, structs são similares a classes em outras linguagens,
// mas sem herança - Go usa composição ao invés de herança.
package models

// User representa um usuário do sistema.
// As tags `json:"..."` definem como o campo será serializado/deserializado em JSON.
// A tag `binding:"required"` é usada pelo Gin para validação automática.
type User struct {
	// ID é o identificador único do usuário.
	// Em Go, campos que começam com letra maiúscula são "exportados" (públicos).
	// Campos com letra minúscula seriam privados ao pacote.
	ID string `json:"id"`

	// Username é o nome de usuário para login.
	// A tag binding:"required" faz o Gin rejeitar requests sem este campo.
	Username string `json:"username" binding:"required"`

	// Email do usuário.
	Email string `json:"email" binding:"required"`

	// Password é a senha do usuário.
	// Note que usamos `json:"-"` para NUNCA retornar a senha em responses JSON.
	// O "-" significa "ignore este campo na serialização JSON".
	Password string `json:"password,omitempty" binding:"required"`

	// PasswordHash é o hash bcrypt da senha.
	// Também não é retornado em JSON por segurança.
	PasswordHash string `json:"-"`
}

// LoginRequest representa os dados necessários para fazer login.
// Criamos uma struct separada para ter validações específicas do login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest representa os dados necessários para registrar um usuário.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse é a struct usada para retornar dados do usuário sem a senha.
// É uma boa prática ter structs separadas para request e response.
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ToResponse converte um User para UserResponse.
// Este é um "método" em Go - uma função associada a uma struct.
// O "(u User)" antes do nome da função é chamado de "receiver".
func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
