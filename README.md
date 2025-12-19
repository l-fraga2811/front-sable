# Go API com Autenticação JWT e CRUD

Uma API REST simples em Go com comentários explicativos para aprendizado.

## 🚀 Funcionalidades

- **Autenticação JWT**: Login, Register e Logout
- **Middleware de Autorização**: Proteção de rotas com JWT
- **CRUD Completo**: Create, Read, Update, Delete de items
- **Armazenamento em Memória**: Simples para aprendizado (dados perdidos ao reiniciar)

## 📁 Estrutura do Projeto

```
api/
├── main.go              # Ponto de entrada e configuração de rotas
├── go.mod               # Gerenciador de dependências
├── handlers/
│   ├── authHandler.go   # Handlers de autenticação (login, register, logout)
│   └── itemHandler.go   # Handlers do CRUD de items
├── middleware/
│   └── authMiddleware.go # Middleware de autenticação JWT
├── models/
│   ├── user.go          # Model de usuário
│   └── item.go          # Model de item
└── storage/
    └── storage.go       # Armazenamento em memória (thread-safe)
```

## 🛠️ Pré-requisitos

- Go 1.21 ou superior

## 📦 Instalação

```bash
# Entre na pasta da API
cd api

# Baixe as dependências
go mod tidy

# Execute a API
go run main.go
```

A API estará disponível em `http://localhost:8080`

## 📚 Endpoints

### Públicos (sem autenticação)

| Método | Endpoint         | Descrição                     |
| ------ | ---------------- | ----------------------------- |
| POST   | `/auth/register` | Registrar novo usuário        |
| POST   | `/auth/login`    | Login e obtenção do token JWT |
| GET    | `/health`        | Health check da API           |

### Protegidos (requerem token JWT)

| Método | Endpoint            | Descrição                 |
| ------ | ------------------- | ------------------------- |
| POST   | `/api/auth/logout`  | Logout (invalida o token) |
| GET    | `/api/auth/profile` | Obter dados do usuário    |
| GET    | `/api/items`        | Listar todos os items     |
| GET    | `/api/items/:id`    | Buscar item por ID        |
| POST   | `/api/items`        | Criar novo item           |
| PUT    | `/api/items/:id`    | Atualizar item            |
| DELETE | `/api/items/:id`    | Deletar item              |

## 🔐 Autenticação

A API usa JWT (JSON Web Tokens) para autenticação. Após o login, você receberá um token que deve ser enviado no header `Authorization` de todas as requisições protegidas.

```
Authorization: Bearer <seu-token-jwt>
```

## 📝 Exemplos de Uso

### Registrar Usuário

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "joao",
    "email": "joao@email.com",
    "password": "senha123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "joao",
    "password": "senha123"
  }'
```

Resposta:

```json
{
  "message": "Login realizado com sucesso",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresAt": "2024-12-18T20:30:00Z",
  "user": {
    "id": "uuid-do-usuario",
    "username": "joao",
    "email": "joao@email.com"
  }
}
```

### Criar Item (autenticado)

```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "title": "Meu primeiro item",
    "description": "Descrição do item",
    "price": 99.90
  }'
```

### Listar Items (autenticado)

```bash
curl http://localhost:8080/api/items \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"
```

### Atualizar Item (autenticado)

```bash
curl -X PUT http://localhost:8080/api/items/ID_DO_ITEM \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "title": "Título atualizado",
    "completed": true
  }'
```

### Deletar Item (autenticado)

```bash
curl -X DELETE http://localhost:8080/api/items/ID_DO_ITEM \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"
```

### Logout (autenticado)

```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"
```

## 📖 Conceitos Go Demonstrados

Este projeto demonstra vários conceitos importantes de Go:

1. **Packages**: Organização de código em pacotes
2. **Structs**: Definição de tipos de dados
3. **Methods**: Funções associadas a structs
4. **Interfaces**: Contratos para tipos
5. **Goroutines e Mutex**: Concorrência segura
6. **Error Handling**: Tratamento de erros idiomático
7. **Defer**: Execução adiada de funções
8. **Maps e Slices**: Estruturas de dados
9. **Pointers**: Referências a valores
10. **Tags de Struct**: Metadados para serialização

## ⚠️ Notas Importantes

- Esta é uma API de **aprendizado**. Para produção, considere:
  - Usar um banco de dados real (PostgreSQL, MySQL, MongoDB)
  - Armazenar a chave JWT em variáveis de ambiente
  - Implementar rate limiting
  - Adicionar validações mais robustas
  - Usar HTTPS
  - Implementar refresh tokens

## 📄 Licença

MIT
