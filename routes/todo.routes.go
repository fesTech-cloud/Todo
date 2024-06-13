package routes

import (
	controller "github.com/festech-cloud/todo/controller"
	"github.com/gin-gonic/gin"
)

func TodoRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/todos", controller.GetTodos())
	incomingRoutes.GET("/todos/:todo_id", controller.GetTodo())
	incomingRoutes.POST("/todos", controller.CreateTodo())
	incomingRoutes.PATCH("/todos/:todo_id", controller.UpdateTodo())
	incomingRoutes.DELETE("/todos/:todo_id", controller.DeleteTodo())
}
