package controllers

import (
	"context"
	// "fmt"
	"go-crud-api/configs"
	"go-crud-api/models"
	"go-crud-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var employeeCollection *mongo.Collection = configs.GetCollection(configs.DB, "test")

func CreateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var employee models.Employee

		// validate that request body is JSON
		if err := c.BindJSON(&employee); err != nil {
			c.JSON(http.StatusBadRequest, responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		// validate that all fields exist (required tag)
		if fieldValidationErr := validate.Struct(&employee); fieldValidationErr != nil {
			c.JSON(http.StatusBadRequest, responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": fieldValidationErr.Error()}})
			return
		}
		newEmployee := models.Employee{
			Id:        primitive.NewObjectID(),
			Firstname: employee.Firstname,
			Lastname:  employee.Lastname,
			Position:  employee.Position,
			Age:       employee.Age}
		result, err := employeeCollection.InsertOne(ctx, newEmployee)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusCreated, responses.EmployeeResponse{Status: http.StatusCreated, Message: "Create Employee Success", Data: map[string]interface{}{"data": result}})
	}
}
func GetAllEmployees() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var allEmployees []models.Employee

		cursor, err := employeeCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		// reading document in batch
		defer cursor.Close(ctx)
		// when the cursor is not exhausted, keeps the iteration
		for cursor.Next(ctx) {
			var employee models.Employee
			if err = cursor.Decode(&employee); err != nil {
				c.JSON(http.StatusInternalServerError, responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			allEmployees = append(allEmployees, employee)
		}
		c.JSON(http.StatusOK, responses.EmployeeResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": allEmployees}})
	}
}
func GetAnEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		employeeId := c.Param("id")
		var employee models.Employee

		// TODO: change this to running number id
		objId, _ := primitive.ObjectIDFromHex(employeeId)
		err := employeeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&employee)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.EmployeeResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": employee}})
	}
}
func DeleteAnEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		employeeId := c.Param("id")

		objId, _ := primitive.ObjectIDFromHex(employeeId)
		result, err := employeeCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusOK, responses.EmployeeResponse{Status: http.StatusOK, Message: "Warning Zero Deletion occured", Data: map[string]interface{}{"data": result}})
			return
		}
		c.JSON(http.StatusOK, responses.EmployeeResponse{Status: http.StatusOK, Message: "Successfully Deleted", Data: map[string]interface{}{"data": map[string]interface{}{"result": result, "deleted": employeeId}}})
	}
}
func UpdateAnEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		employeeId := c.Param("id")
		var employee models.Employee

		objId, _ := primitive.ObjectIDFromHex(employeeId)

		// Validate that request is JSON
		if err := c.BindJSON(&employee); err != nil {
			c.JSON(http.StatusBadRequest, responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		// Validate required fields
		if validationErr := validate.Struct(&employee); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.EmployeeResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		updateData := bson.M{"firstname": employee.Firstname, "lastname": employee.Lastname, "position": employee.Position, "age": employee.Age}

		result, err := employeeCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": updateData})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		var updatedEmployee models.Employee
		// if there is matched
		if result.MatchedCount == 1 {
			err = employeeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedEmployee)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.EmployeeResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}
		c.JSON(http.StatusOK, responses.EmployeeResponse{Status: http.StatusOK, Message: "Update Success", Data: map[string]interface{}{"data": updatedEmployee}})
	}
}
