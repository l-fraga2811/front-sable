# Items Manager - Frontend

Interface web para gerenciamento de items, construída com NextJS 16 e integrada com API Go.

## Stack Tecnológica

- **NextJS 16** - Framework React com App Router
- **TailwindCSS 4** - Estilização
- **shadcn/ui** - Componentes UI
- **Redux Toolkit** - Gerenciamento de estado
- **Axios** - Cliente HTTP
- **Lucide React** - Ícones

## Pré-requisitos

- Node.js 18+
- API Go rodando em `http://localhost:8080`

## Instalação

```bash
npm install
```

## Configuração

Crie um arquivo `.env.local` na raiz do projeto:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Executando

```bash
npm run dev
```

Acesse [http://localhost:3000](http://localhost:3000)

## Estrutura do Projeto

```
src/
├── app/                    # Páginas (App Router)
│   ├── login/             # Página de login
│   ├── register/          # Página de registro
│   └── dashboard/         # Área autenticada
│       └── profile/       # Página de perfil
├── components/
│   ├── ui/                # Componentes shadcn/ui
│   └── shared/            # Componentes compartilhados
├── lib/
│   └── axios.ts           # Configuração do Axios
├── store/                 # Redux Store
│   ├── auth/              # Módulo de autenticação
│   └── items/             # Módulo de items
└── types/                 # Tipos TypeScript
```

## Funcionalidades

- **Autenticação** - Login, registro e logout com JWT
- **CRUD de Items** - Criar, listar, editar e excluir items
- **Perfil** - Visualização de dados do usuário
- **Proteção de Rotas** - Redirecionamento automático para login
