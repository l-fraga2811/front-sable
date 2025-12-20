package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RestClient é o cliente HTTP para a API REST do Supabase (PostgREST).
//
// CONCEITO: PostgREST
// O Supabase expõe seu banco PostgreSQL via uma API REST automática chamada
// PostgREST. Cada tabela vira um endpoint:
//   - GET /rest/v1/items → SELECT * FROM items
//   - POST /rest/v1/items → INSERT INTO items
//   - PATCH /rest/v1/items?id=eq.123 → UPDATE items WHERE id = '123'
//   - DELETE /rest/v1/items?id=eq.123 → DELETE FROM items WHERE id = '123'
//
// AUTENTICAÇÃO:
//   - apikey: Chave anon do projeto (sempre enviada)
//   - Authorization: Bearer <access_token> (token do usuário logado)
//
// RLS (Row Level Security):
// O Supabase usa RLS para controlar acesso. O token do usuário determina
// quais linhas ele pode ver/modificar baseado nas policies definidas no banco.
type RestClient struct {
	projectURL string
	anonKey    string
	client     *http.Client
}

// NewRestClient cria um novo cliente REST.
func NewRestClient(cfg Config) *RestClient {
	return &RestClient{
		projectURL: strings.TrimRight(cfg.ProjectURL, "/"),
		anonKey:    cfg.AnonKey,
		client:     &http.Client{Timeout: 15 * time.Second},
	}
}

// ItemRow representa uma linha da tabela "items" no Supabase.
//
// CONVENÇÃO: snake_case no banco, camelCase no Go
// O PostgreSQL usa snake_case (user_id, created_at).
// As tags json mapeiam para os nomes das colunas no banco.
type ItemRow struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       Numeric `json:"price"`
	Completed   bool    `json:"completed"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// CreateItemPayload é o payload para criar um novo item.
type CreateItemPayload struct {
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Completed   bool    `json:"completed"`
}

// UpdateItemPayload é o payload para atualizar um item.
//
// PADRÃO: Ponteiros para campos opcionais
// Usamos ponteiros (*string, *float64) para diferenciar entre:
//   - Campo não enviado (nil) → não atualiza
//   - Campo enviado com valor zero ("", 0) → atualiza para zero
//
// A tag `omitempty` omite campos nil do JSON.
type UpdateItemPayload struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Completed   *bool    `json:"completed,omitempty"`
}

// ListItems retorna todos os items do usuário (filtrado por RLS).
//
// QUERY POSTGREST:
// GET /rest/v1/items?select=*&order=created_at.desc
//
// O RLS garante que só retorna items do usuário autenticado.
func (c *RestClient) ListItems(ctx context.Context, accessToken string) ([]ItemRow, error) {
	q := url.Values{}
	q.Set("select", "*")
	q.Set("order", "created_at.desc")
	return c.getItems(ctx, accessToken, q)
}

// GetItemByID busca um item específico pelo ID.
//
// QUERY POSTGREST:
// GET /rest/v1/items?select=*&id=eq.{id}
//
// O operador "eq." significa "equals" (igual a).
func (c *RestClient) GetItemByID(ctx context.Context, accessToken string, id string) (ItemRow, bool, error) {
	q := url.Values{}
	q.Set("select", "*")
	q.Set("id", "eq."+id)

	items, err := c.getItems(ctx, accessToken, q)
	if err != nil {
		return ItemRow{}, false, err
	}
	if len(items) == 0 {
		return ItemRow{}, false, nil
	}
	return items[0], true, nil
}

// CreateItem insere um novo item no banco.
//
// QUERY POSTGREST:
// POST /rest/v1/items
// Body: {"user_id": "...", "title": "...", ...}
// Header: Prefer: return=representation (retorna o item criado)
func (c *RestClient) CreateItem(ctx context.Context, accessToken string, payload CreateItemPayload) (ItemRow, error) {
	var created []ItemRow
	if err := c.doJSON(ctx, http.MethodPost, "/rest/v1/items", accessToken, nil, payload, &created, map[string]string{"Prefer": "return=representation"}); err != nil {
		return ItemRow{}, err
	}
	if len(created) == 0 {
		return ItemRow{}, errors.New("supabase não retornou item criado")
	}
	return created[0], nil
}

// UpdateItem atualiza um item existente.
//
// QUERY POSTGREST:
// PATCH /rest/v1/items?id=eq.{id}
// Body: {"title": "novo titulo", ...}
func (c *RestClient) UpdateItem(ctx context.Context, accessToken string, id string, payload UpdateItemPayload) (ItemRow, bool, error) {
	q := url.Values{}
	q.Set("id", "eq."+id)

	var updated []ItemRow
	err := c.doJSON(ctx, http.MethodPatch, "/rest/v1/items", accessToken, q, payload, &updated, map[string]string{"Prefer": "return=representation"})
	if err != nil {
		return ItemRow{}, false, err
	}
	if len(updated) == 0 {
		return ItemRow{}, false, nil
	}
	return updated[0], true, nil
}

// DeleteItem remove um item do banco.
//
// QUERY POSTGREST:
// DELETE /rest/v1/items?id=eq.{id}
func (c *RestClient) DeleteItem(ctx context.Context, accessToken string, id string) (bool, error) {
	q := url.Values{}
	q.Set("id", "eq."+id)

	var deleted []ItemRow
	if err := c.doJSON(ctx, http.MethodDelete, "/rest/v1/items", accessToken, q, nil, &deleted, map[string]string{"Prefer": "return=representation"}); err != nil {
		return false, err
	}
	return len(deleted) > 0, nil
}

// getItems é um helper interno para buscar items com query customizada.
func (c *RestClient) getItems(ctx context.Context, accessToken string, q url.Values) ([]ItemRow, error) {
	var out []ItemRow
	if err := c.doJSON(ctx, http.MethodGet, "/rest/v1/items", accessToken, q, nil, &out, nil); err != nil {
		return nil, err
	}
	if out == nil {
		out = []ItemRow{}
	}
	return out, nil
}

// doJSON faz uma requisição HTTP e deserializa a resposta JSON.
//
// PARÂMETROS:
//   - ctx: Contexto para cancelamento/timeout
//   - method: GET, POST, PATCH, DELETE
//   - path: Caminho da API (ex: /rest/v1/items)
//   - accessToken: Token JWT do usuário
//   - q: Query parameters (filtros)
//   - payload: Body da requisição (será serializado para JSON)
//   - out: Ponteiro para deserializar a resposta
//   - extraHeaders: Headers adicionais (ex: Prefer)
func (c *RestClient) doJSON(ctx context.Context, method string, path string, accessToken string, q url.Values, payload any, out any, extraHeaders map[string]string) error {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	resp, err := c.do(ctx, method, path, accessToken, q, body, extraHeaders)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		if len(b) == 0 {
			return errors.New("requisição supabase falhou")
		}
		return errors.New(string(b))
	}

	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// do é o método base que monta e executa a requisição HTTP.
//
// HEADERS ENVIADOS:
//   - apikey: Chave anon do projeto (obrigatório para Supabase)
//   - Authorization: Bearer token (para RLS)
//   - Content-Type: application/json (se tiver body)
func (c *RestClient) do(ctx context.Context, method string, path string, accessToken string, q url.Values, body io.Reader, extraHeaders map[string]string) (*http.Response, error) {
	base := strings.TrimRight(c.projectURL, "/")
	fullURL := base + path
	if q != nil {
		fullURL = fullURL + "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.anonKey)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("resposta vazia")
	}
	return resp, nil
}
