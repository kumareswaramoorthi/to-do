package router

import (
	"todo/controller"
	"todo/database"
	"todo/middleware"
	"todo/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	// Get Db connection
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	//create access variables
	todoDatabase := database.NewTodoDatabase(db)
	todoService := services.NewTodoService(todoDatabase)
	todoController := controller.NewTodoController(todoService)
	todo := router.Group("/api/todo/v1/")
	{
		todo.POST("/signup", todoController.SignUp)
		todo.POST("/signin", todoController.SignIn)
		todo.POST("/addtodo", middleware.TokenAuthMiddleware(), todoController.AddTodo)
		todo.DELETE("/deletetodo", middleware.TokenAuthMiddleware(), todoController.DeleteTodo)
		todo.PUT("/edittodo", middleware.TokenAuthMiddleware(), todoController.EditTodo)
		todo.GET("/getalltodos", middleware.TokenAuthMiddleware(), todoController.GetAllTodos)
		todo.GET("/gettodobycategory", middleware.TokenAuthMiddleware(), todoController.GetTodoByCategory)
		todo.POST("/marktodo", middleware.TokenAuthMiddleware(), todoController.MarkTodo)
		todo.POST("/addcategory", middleware.TokenAuthMiddleware(), todoController.AddCategory)
		todo.GET("/getcategory", middleware.TokenAuthMiddleware(), todoController.GetCategory)
		todo.DELETE("/deletecategory", middleware.TokenAuthMiddleware(), todoController.DeleteCategory)
	}
	return router
}
