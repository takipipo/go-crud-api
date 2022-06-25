package main

import (
	"go-crud-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.EmployeeRoute(router)
	router.Run("localhost:8000")
}