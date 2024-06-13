package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/festech-cloud/todo/database"
	"github.com/festech-cloud/todo/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todoCollection *mongo.Collection
var validate = validator.New()

func Init(client *mongo.Client) {
	todoCollection = database.OpenCollection(client, "Todo")
}

func GetTodos() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		result, err := todoCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": false, "message": "something happened while fetching todos"})
			return
		}

		var allTodos []bson.M
		err = result.All(ctx, &allTodos)

		if err != nil {
			log.Fatal(err)
		}

		defer cancel()
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Todo fetched successfully", "data": allTodos})

		// c.JSON(http.StatusOK, gin.H{"status": true, "message": "route working successfully"})
	}
}

func GetTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var todo models.Todo
		todo_Id := c.Param("todo_id")

		err := todoCollection.FindOne(ctx, bson.M{"todo_id": todo_Id}).Decode(&todo)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": false, "message": "An error occured while fetching todo"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true, "message": "Todo fetched successfully", "data": todo})

	}
}

func CreateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("App is runing ")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var todo models.Todo
		err := c.BindJSON(&todo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
		}

		validatorErr := validate.Struct(todo)
		if validatorErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validatorErr.Error()})
			return
		}

		todo.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		todo.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		todo.ID = primitive.NewObjectID()
		todo.Todo_id = todo.ID.Hex()

		result, insertErr := todoCollection.InsertOne(ctx, todo)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "An error occured while creating menu"})
			return
		}
		c.JSON(http.StatusOK, result)
		defer cancel()

	}
}

func UpdateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		todo_id := c.Param("todo_id")
		var todo models.Todo
		var filter = bson.M{"todo_id": todo_id}
		var updateTodo primitive.D
		err := c.BindJSON(&todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "An error occured, pls try again"})
			return
		}
		err = validate.Struct(&todo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		if todo.Todo != "" {
			updateTodo = append(updateTodo, bson.E{"todo", todo.Todo})
		}
		if todo.Completed != nil {
			updateTodo = append(updateTodo, bson.E{"completed", todo.Completed})
		}
		todo.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		updateTodo = append(updateTodo, bson.E{"updated_at", todo.Updated_at})
		_, err = todoCollection.UpdateOne(ctx, filter, bson.D{
			{"$set", updateTodo},
		},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Could not update todo"})
			defer cancel()
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"status": true, "message": "update successful"})
	}
}

func DeleteTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		todoId := c.Param("todo_id")
		filter := bson.M{"todo_id": todoId}

		_, err := todoCollection.DeleteOne(ctx, filter)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": false, "message": "An error occured while deleting todo"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": true,
			"message": "Deleted successfully"})
	}
}
