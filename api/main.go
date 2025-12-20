package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-api/docs"
	"go-api/middleware"
	"go-api/modules/auth"

	itemController "go-api/modules/item/controller"
	itemRoutes "go-api/modules/item/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis do sistema")
	}
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()
	router.Use(corsMiddleware())

	authHandler := handlers.NewAuthHandler()
	itemCtrl := itemController.NewItemController()

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.GET("/auth/profile", authHandler.GetProfile)
		itemRoutes.Register(protectedRoutes, itemCtrl)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API está funcionando!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciando na porta %s", port)
	log.Printf("Documentação: http://localhost:%s/docs/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
