# Criação da Interface NextJS para API Go

## Objetivo

Criar uma interface web moderna com NextJS para consumir a API Go existente, implementando autenticação JWT e CRUD de items.

## Análise da API

### Endpoints Públicos

- `POST /auth/register` - Registrar usuário (username, email, password)
- `POST /auth/login` - Login (username, password) → retorna token JWT

### Endpoints Protegidos (requerem Bearer Token)

- `POST /api/auth/logout` - Logout
- `GET /api/auth/profile` - Obter perfil do usuário
- `GET /api/items` - Listar items do usuário
- `GET /api/items/:id` - Buscar item por ID
- `POST /api/items` - Criar item (title, description, price)
- `PUT /api/items/:id` - Atualizar item (title, description, price, completed)
- `DELETE /api/items/:id` - Deletar item

### Modelos de Dados

**User:**

- id: string
- username: string
- email: string

**Item:**

- id: string
- title: string
- description: string
- price: number
- completed: boolean
- createdAt: Date
- updatedAt: Date
- userId: string

## Stack Tecnológica

- **Framework:** NextJS 14+ (App Router)
- **Estilização:** TailwindCSS
- **Componentes:** shadcn/ui
- **Ícones:** Lucide React
- **Estado Global:** Redux Toolkit
- **HTTP Client:** Axios

## Estrutura do Projeto

```
web/
├── src/
│   ├── app/
│   │   ├── (auth)/
│   │   │   ├── login/
│   │   │   └── register/
│   │   ├── (dashboard)/
│   │   │   ├── items/
│   │   │   └── profile/
│   │   ├── layout.tsx
│   │   └── page.tsx
│   ├── components/
│   │   ├── ui/
│   │   └── shared/
│   ├── lib/
│   │   └── axios.ts
│   ├── store/
│   │   ├── index.ts
│   │   ├── auth/
│   │   │   ├── actions.ts
│   │   │   ├── reducers.ts
│   │   │   ├── services.ts
│   │   │   └── selectors.ts
│   │   └── items/
│   │       ├── actions.ts
│   │       ├── reducers.ts
│   │       ├── services.ts
│   │       └── selectors.ts
│   └── types/
│       └── index.ts
```

## Checklist de Implementação

- [x] 1. Criar projeto NextJS com dependências
- [x] 2. Configurar TailwindCSS e shadcn/ui
- [x] 3. Criar tipos TypeScript
- [x] 4. Configurar Axios com interceptors
- [x] 5. Configurar Redux Store
- [x] 6. Implementar módulo Auth (actions, reducers, services, selectors)
- [x] 7. Implementar módulo Items (actions, reducers, services, selectors)
- [x] 8. Criar componentes UI base
- [x] 9. Criar página de Login
- [x] 10. Criar página de Registro
- [x] 11. Criar página de Dashboard/Items
- [x] 12. Criar página de Perfil
- [x] 13. Implementar proteção de rotas
- [x] 14. Testar integração completa
