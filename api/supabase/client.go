package supabase

import "sync"

// Client é o cliente principal do Supabase que agrupa todos os componentes.
//
// COMPONENTES:
//   - Config: Configurações carregadas do ambiente
//   - TokenValidator: Valida tokens JWT do Supabase Auth
//   - RestClient: Faz requisições à API REST (PostgREST) do Supabase
type Client struct {
	Config         Config
	TokenValidator *TokenValidator
	RestClient     *RestClient
}

// Variáveis para implementar o padrão Singleton.
//
// PADRÃO: Singleton com sync.Once
// sync.Once garante que a inicialização aconteça apenas UMA vez,
// mesmo com múltiplas goroutines chamando GetClient() simultaneamente.
// Isso é importante porque:
//   - Evita criar múltiplas conexões desnecessárias
//   - Garante consistência de configuração
//   - É thread-safe (seguro para uso concorrente)
var (
	clientOnce sync.Once
	clientInst *Client
	clientErr  error
)

// GetClient retorna a instância única do cliente Supabase.
//
// PADRÃO: Lazy Initialization
// O cliente só é criado na primeira chamada, não no início do programa.
// Isso é útil porque:
//   - Evita erros de inicialização se o Supabase não for usado
//   - Permite que as variáveis de ambiente sejam carregadas antes
//
// RETORNO: (*Client, error)
// Em Go, é comum retornar (valor, erro). O chamador deve sempre verificar
// o erro antes de usar o valor.
func GetClient() (*Client, error) {
	clientOnce.Do(func() {
		cfg, err := LoadConfig()
		if err != nil {
			clientErr = err
			return
		}

		clientInst = &Client{
			Config:         cfg,
			TokenValidator: NewTokenValidator(cfg),
			RestClient:     NewRestClient(cfg),
		}
	})
	return clientInst, clientErr
}
