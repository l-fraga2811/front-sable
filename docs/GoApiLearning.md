# Plano de Ação - API Go com Autenticação e CRUD

## Objetivo

Criar uma API REST em Go com comentários explicativos para aprendizado, incluindo:

- Sistema de autenticação (login, register, logout)
- Middleware de autorização JWT
- CRUD completo de um recurso

## Stack Tecnológica

- **Go** - Linguagem principal
- **Gin** - Framework web (leve e rápido)
- **JWT** - Autenticação via tokens
- **bcrypt** - Hash de senhas
- **In-memory storage** - Armazenamento simples para aprendizado

## Estrutura de Pastas

```
/api
├── main.go              # Ponto de entrada da aplicação
├── go.mod               # Gerenciador de dependências
├── handlers/
│   ├── authHandler.go   # Handlers de autenticação
│   └── itemHandler.go   # Handlers do CRUD
├── middleware/
│   └── authMiddleware.go # Middleware JWT
├── models/
│   ├── user.go          # Model de usuário
│   └── item.go          # Model do CRUD
├── storage/
│   └── storage.go       # Armazenamento em memória
└── README.md            # Instruções de uso
```

## Checklist

- [x] Criar estrutura de pastas do projeto Go
- [x] Criar arquivo go.mod e instalar dependências
- [x] Criar models (User e Item para CRUD)
- [x] Criar storage em memória
- [x] Criar middleware de autenticação JWT
- [x] Criar handlers de autenticação (login, register, logout)
- [x] Criar handlers CRUD (Create, Read, Update, Delete)
- [x] Criar arquivo main.go com rotas
- [x] Criar README com instruções

## Endpoints Planejados

### Autenticação (públicos)

- `POST /auth/register` - Registrar novo usuário
- `POST /auth/login` - Login e retorno do token JWT
- `POST /auth/logout` - Logout (invalidar token)

### CRUD Items (protegidos por JWT)

- `GET /api/items` - Listar todos os items
- `GET /api/items/:id` - Buscar item por ID
- `POST /api/items` - Criar novo item
- `PUT /api/items/:id` - Atualizar item
- `DELETE /api/items/:id` - Deletar item
