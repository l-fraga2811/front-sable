// Package storage fornece armazenamento em memória para a aplicação.
// Em produção, você substituiria isso por um banco de dados real (PostgreSQL, MySQL, MongoDB, etc).
// Este storage usa sync.RWMutex para ser thread-safe (seguro para uso concorrente).
package storage

import (
	"errors"
	"sync"

	"go-api/models"
)

// Erros personalizados.
// Em Go, é comum criar variáveis de erro reutilizáveis.
// O prefixo "Err" é uma convenção para variáveis de erro.
var (
	ErrUserNotFound     = errors.New("usuário não encontrado")
	ErrUserExists       = errors.New("usuário já existe")
	ErrItemNotFound     = errors.New("item não encontrado")
	ErrInvalidToken     = errors.New("token inválido")
	ErrUnauthorized     = errors.New("não autorizado")
)

// Storage é a estrutura principal que mantém todos os dados em memória.
// sync.RWMutex permite múltiplas leituras simultâneas, mas apenas uma escrita por vez.
type Storage struct {
	// mu é o mutex para controle de concorrência.
	// RWMutex = Read-Write Mutex
	// - RLock/RUnlock: permite múltiplas goroutines lerem ao mesmo tempo
	// - Lock/Unlock: bloqueia todas as outras goroutines (leitura e escrita)
	mu sync.RWMutex

	// users armazena os usuários. A chave é o username.
	// map[string]models.User significa: mapa de string para User
	users map[string]models.User

	// items armazena os items. A chave é o ID do item.
	items map[string]models.Item

	// blacklistedTokens armazena tokens que foram invalidados (logout).
	// Usamos map[string]bool como um "Set" (conjunto).
	blacklistedTokens map[string]bool
}

// store é a instância global do storage.
// Em Go, variáveis em minúsculo são privadas ao pacote.
var store *Storage

// once garante que a inicialização aconteça apenas uma vez.
// sync.Once é útil para inicialização lazy e thread-safe.
var once sync.Once

// GetStorage retorna a instância única do storage (Singleton pattern).
// O padrão Singleton garante que existe apenas uma instância do storage.
func GetStorage() *Storage {
	// once.Do executa a função apenas na primeira chamada.
	// Chamadas subsequentes não executam a função novamente.
	once.Do(func() {
		store = &Storage{
			users:             make(map[string]models.User),
			items:             make(map[string]models.Item),
			blacklistedTokens: make(map[string]bool),
		}
	})
	return store
}

// CreateUser adiciona um novo usuário ao storage.
// Retorna erro se o usuário já existir.
func (s *Storage) CreateUser(user models.User) error {
	// Lock() bloqueia para escrita - nenhuma outra goroutine pode ler ou escrever.
	s.mu.Lock()
	// defer garante que Unlock() será chamado quando a função terminar,
	// mesmo se houver um panic. É como um "finally" em outras linguagens.
	defer s.mu.Unlock()

	// Verifica se o usuário já existe.
	// O segundo valor retornado (exists) indica se a chave existe no map.
	if _, exists := s.users[user.Username]; exists {
		return ErrUserExists
	}

	s.users[user.Username] = user
	return nil
}

// GetUserByUsername busca um usuário pelo username.
// Retorna o usuário e um erro (nil se encontrado).
func (s *Storage) GetUserByUsername(username string) (models.User, error) {
	// RLock() permite múltiplas leituras simultâneas.
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		// Retornamos um User vazio e o erro.
		// Em Go, é comum retornar o "zero value" junto com um erro.
		return models.User{}, ErrUserNotFound
	}
	return user, nil
}

// CreateItem adiciona um novo item ao storage.
func (s *Storage) CreateItem(item models.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[item.ID] = item
	return nil
}

// GetItemByID busca um item pelo ID.
func (s *Storage) GetItemByID(id string) (models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return models.Item{}, ErrItemNotFound
	}
	return item, nil
}

// GetItemsByUserID retorna todos os items de um usuário específico.
func (s *Storage) GetItemsByUserID(userID string) []models.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Criamos um slice vazio para armazenar os resultados.
	// []models.Item{} é um slice literal vazio.
	var items []models.Item

	// range itera sobre o map. _ ignora a chave (ID), item é o valor.
	for _, item := range s.items {
		if item.UserID == userID {
			// append adiciona um elemento ao slice.
			items = append(items, item)
		}
	}

	return items
}

// UpdateItem atualiza um item existente.
func (s *Storage) UpdateItem(item models.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[item.ID]; !exists {
		return ErrItemNotFound
	}

	s.items[item.ID] = item
	return nil
}

// DeleteItem remove um item pelo ID.
func (s *Storage) DeleteItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return ErrItemNotFound
	}

	// delete é uma função built-in para remover chaves de um map.
	delete(s.items, id)
	return nil
}

// BlacklistToken adiciona um token à lista negra (logout).
func (s *Storage) BlacklistToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.blacklistedTokens[token] = true
}

// IsTokenBlacklisted verifica se um token está na lista negra.
func (s *Storage) IsTokenBlacklisted(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.blacklistedTokens[token]
}
