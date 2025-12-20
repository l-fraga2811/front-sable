package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-api/docs" 
	"go-api/handlers"
	"go-api/middleware"
)

func main() {
	// Carrega variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis do sistema")
	}

	// Configura o modo do Gin.
	// gin.ReleaseMode desativa logs de debug em produção.
	// gin.DebugMode (padrão) mostra logs detalhados.
	// Você pode definir via variável de ambiente: GIN_MODE=release
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// gin.Default() cria um router com middlewares padrão:
	// - Logger: loga todas as requisições
	// - Recovery: recupera de panics e retorna 500
	router := gin.Default()

	// Configura CORS (Cross-Origin Resource Sharing).
	// Isso permite que frontends em outros domínios acessem a API.
	router.Use(corsMiddleware())

	// Cria as instâncias dos handlers.
	// Usamos o padrão de injeção de dependência.
	authHandler := handlers.NewAuthHandler()
	itemHandler := handlers.NewItemHandler()

	// ============================================
	// ROTAS PÚBLICAS (não requerem autenticação)
	// ============================================

	// ============================================
	// ROTAS PROTEGIDAS (requerem autenticação)
	// ============================================

	// Grupo de rotas que requerem autenticação.
	// middleware.AuthMiddleware() é executado ANTES de cada handler neste grupo.
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		// GET /auth/profile - Obter dados do usuário autenticado
		protectedRoutes.GET("/auth/profile", authHandler.GetProfile)

		// ============================================
		// CRUD DE ITEMS
		// ============================================

		// GET /api/items - Listar todos os items do usuário
		protectedRoutes.GET("/items", itemHandler.GetAll)

		// GET /api/items/:id - Buscar item por ID
		// :id é um parâmetro de rota (path parameter)
		protectedRoutes.GET("/items/:id", itemHandler.GetByID)

		// POST /api/items - Criar novo item
		protectedRoutes.POST("/items", itemHandler.Create)

		// PUT /api/items/:id - Atualizar item existente
		protectedRoutes.PUT("/items/:id", itemHandler.Update)

		// DELETE /api/items/:id - Deletar item
		protectedRoutes.DELETE("/items/:id", itemHandler.Delete)
	}

	// Rota de documentação Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rota de health check (útil para monitoramento).
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API está funcionando!",
		})
	})

	// Define a porta do servidor.
	// Primeiro tenta ler da variável de ambiente PORT.
	// Se não existir, usa 8080 como padrão.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Loga a mensagem de início.
	log.Printf("🚀 Servidor iniciando na porta %s", port)
	log.Printf("📚 Documentação: http://localhost:%s/health", port)

	// router.Run inicia o servidor HTTP.
	// Ele bloqueia a execução até o servidor ser encerrado.
	// O formato ":8080" significa "escute em todas as interfaces na porta 8080".
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// corsMiddleware configura o CORS para a API.
// CORS é necessário quando o frontend está em um domínio diferente da API.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Access-Control-Allow-Origin define quais origens podem acessar a API.
		// "*" permite qualquer origem (não recomendado em produção).
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Access-Control-Allow-Credentials permite envio de cookies/auth headers.
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Access-Control-Allow-Headers define quais headers o cliente pode enviar.
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		// Access-Control-Allow-Methods define quais métodos HTTP são permitidos.
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// Requisições OPTIONS são "preflight requests" do CORS.
		// O navegador envia OPTIONS antes de requisições "complexas".
		if c.Request.Method == "OPTIONS" {
			// Retorna 204 (No Content) para preflight requests.
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
