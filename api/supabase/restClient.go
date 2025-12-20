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
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateItemPayload struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Completed   *bool    `json:"completed,omitempty"`
}

func (c *RestClient) ListItems(ctx context.Context, accessToken string) ([]ItemRow, error) {
	q := url.Values{}
	q.Set("select", "*")
	q.Set("order", "created_at.desc")
	return c.getItems(ctx, accessToken, q)
}

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

func (c *RestClient) DeleteItem(ctx context.Context, accessToken string, id string) (bool, error) {
	q := url.Values{}
	q.Set("id", "eq."+id)

	var deleted []ItemRow
	if err := c.doJSON(ctx, http.MethodDelete, "/rest/v1/items", accessToken, q, nil, &deleted, map[string]string{"Prefer": "return=representation"}); err != nil {
		return false, err
	}
	return len(deleted) > 0, nil
}

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
