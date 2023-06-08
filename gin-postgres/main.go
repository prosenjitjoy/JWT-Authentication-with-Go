package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"main/routes"
	"net/http"
)

func main() {

	router := gin.Default()
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api/v1", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"success": "Access granted for api/v1"})
	})
	router.GET("/api/v2", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"success": "Access granted for api/v2"})
	})

	if err := router.Run(":5000"); err != nil {
		log.Fatal("Error initializing server", err)
	}
}
