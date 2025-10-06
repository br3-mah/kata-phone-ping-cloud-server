package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize database connection
	if err := initDatabase(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer closeDatabase()

	// Load HTML templates
	loadTemplates()

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/device-ping", handleDevicePing)
		api.GET("/devices", getDevices)
		api.GET("/device/:uuid", getDeviceByUUID)
		api.DELETE("/device/:uuid", deleteDevice)
	}

	// Web interface routes
	r.GET("/", serveIndex)
	r.Static("/static", "./static")

	log.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}
