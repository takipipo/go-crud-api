package routes

import (
	"go-crud-api/controllers"

	"github.com/gin-gonic/gin"
)



func EmployeeRoute(router *gin.Engine) {
	router.GET("/employees", controllers.GetAllEmployees())
	router.GET("/employee/:id", controllers.GetAnEmployee())
	router.POST("/employee", controllers.CreateEmployee())
	router.DELETE("/delete/:id", controllers.DeleteAnEmployee())
	router.PUT("/update/:id", controllers.UpdateAnEmployee())
}