package main

import (
	"log"

	"todo-api/handlers"
	"todo-api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	models.InitDB()

	v1 := r.Group("/api/v1")

	tasks := v1.Group("/tasks")
	{
		tasks.GET("/", handlers.GetTasks)
		tasks.GET("/:id", handlers.GetTaskByID)
		tasks.POST("/", handlers.ValidateToken, handlers.CreateTask)
		tasks.PUT("/:id", handlers.ValidateToken, handlers.UpdateTask)
		tasks.DELETE("/:id", handlers.ValidateToken, handlers.DeleteTask)
	}
	users := v1.Group("/users")
	{
		users.GET("/", handlers.ValidateToken, handlers.GetUsers)
		users.GET("/:id", handlers.GetUserByID)
		users.GET("/:id/tasks", handlers.GetTasksByUserID)
		users.GET("/:id/tasks/overdue", handlers.GetOverdueTasksByUserID)
		users.PUT("/:id", handlers.ValidateToken, handlers.UpdateUser)
		users.DELETE("/:id", handlers.ValidateToken, handlers.DeleteUser)
	}
	auth := r.Group("/auth")
	{
		auth.POST("/sign-in", handlers.SignIn)
		auth.POST("/sign-up", handlers.SignUp)
	}

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
