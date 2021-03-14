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
	ctrl := controller.NewTodoController(todoService)
	todo := router.Group("/api/todo/v1/")
	{
		todo.POST("/signup", ctrl.SignUpController)
		todo.POST("/signin", ctrl.SignInController)
		todo.POST("/addtodo", middleware.TokenAuthMiddleware(), ctrl.AddTodoController)
		todo.DELETE("/deletetodo", middleware.TokenAuthMiddleware(), ctrl.DeleteTodoController)
		todo.PUT("/edittodo", middleware.TokenAuthMiddleware(), ctrl.EditTodoController)
		todo.GET("/getalltodos", middleware.TokenAuthMiddleware(), ctrl.GetAllTodosController)
		todo.GET("/gettodobycategory", middleware.TokenAuthMiddleware(), ctrl.GetTodoByCategoryController)
		todo.POST("/marktodo", middleware.TokenAuthMiddleware(), ctrl.MarkTodoController)
		todo.POST("/addcategory", middleware.TokenAuthMiddleware(), ctrl.AddCategoryController)
		todo.GET("/getcategory", middleware.TokenAuthMiddleware(), ctrl.GetCategoryController)
		todo.DELETE("/deletecategory", middleware.TokenAuthMiddleware(), ctrl.DeleteCategoryController)
	}
	return router
}
