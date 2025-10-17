// main.go
// Punto de entrada de la API REST.
// - Inicializa la conexión a la base de datos
// - Configura el router HTTP (Gin)
// - Registra los handlers y rutas
package main

import (
	"log"
	"mi-api-rest/database"
	"mi-api-rest/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar base de datos
	if err := database.InitDB(); err != nil {
		log.Fatal("Error inicializando la base de datos:", err)
	}

	// Configurar router
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		// Pasar al siguiente middleware
		c.Next()
	})

	// Inicializar handlers
	itemHandler := handlers.NewItemHandler(database.DB)

	// Rutas de la API
	// - POST /api/v1/items : crear un registro
	// - GET  /api/v1/items : listar todos los registros
	// - GET  /api/v1/items/:id : obtener por id
	// - PUT  /api/v1/items : sincronizar/actualizar desde API externa
	api := router.Group("/api/v1")
	{
		items := api.Group("/items")
		{
			items.POST("", itemHandler.CreateItem)
			items.GET("", itemHandler.GetItems)
			items.GET("/:id", itemHandler.GetItem)
			items.PUT("/", itemHandler.UpdateItems)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Error iniciando el servidor:", err)
	}
}
