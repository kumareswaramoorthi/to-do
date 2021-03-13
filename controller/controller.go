package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"todo/model"

	"todo/services"

	"github.com/gin-gonic/gin"
)

var _ TodoCtrl = &todoCtrl{}

type TodoCtrl interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	AddTodo(ctx *gin.Context)
	DeleteTodo(ctx *gin.Context)
	GetAllTodos(ctx *gin.Context)
	MarkTodo(ctx *gin.Context)
	AddCategory(ctx *gin.Context)
	EditTodo(ctx *gin.Context)
	GetTodoByCategory(ctx *gin.Context)
	GetCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
}

type todoCtrl struct {
	todoSrv services.TodoService
}

func NewTodoController(todosrv services.TodoService) TodoCtrl {
	return todoCtrl{
		todoSrv: todosrv,
	}
}

//SignUp controller to create a new user
func (t todoCtrl) SignUp(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	errSignup := t.todoSrv.SignUp(ctx, &user)
	if errSignup != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(errSignup))
		return
	}
	ctx.JSON(http.StatusOK, "Successfully created account")
}

//SignIn controller provides token and user information on successfull login
func (t todoCtrl) SignIn(ctx *gin.Context) {
	var user model.UserLogin
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	response, errSignup := t.todoSrv.SignIn(ctx, &user)
	if errSignup != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "server error")
		return
	}
	ctx.JSON(http.StatusOK, response)
}

//AddTodo controller to add a new todo for a user
func (t todoCtrl) AddTodo(ctx *gin.Context) {
	var todo model.Todo
	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	err := t.todoSrv.AddTodo(ctx, &todo)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	ctx.JSON(http.StatusOK, "Todo Added Successfully")
}

//DeleteTodo controller to delete a todo
func (t todoCtrl) DeleteTodo(ctx *gin.Context) {
	todo := ctx.Query("id")
	errDelete := t.todoSrv.DeleteTodo(ctx, todo)
	if errDelete != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(errDelete))
		return
	}
	ctx.JSON(http.StatusOK, "Todo Deleted Successfully")
}

//GetAllTodos controller to get all todo for a user
func (t todoCtrl) GetAllTodos(ctx *gin.Context) {
	todos, err := t.todoSrv.GetAlltodo(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	ctx.JSON(http.StatusOK, todos)
}

//MarkTodo controller to mark a todo status
func (t todoCtrl) MarkTodo(ctx *gin.Context) {
	var todo model.MarkTodo
	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	errMark := t.todoSrv.MarkTodo(ctx, todo)
	if errMark != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(errMark))
		return
	}
	ctx.JSON(http.StatusOK, "Todo Marked Successfully")
}

//AddCategory controller to add category
func (t todoCtrl) AddCategory(ctx *gin.Context) {
	var category model.Category
	if err := ctx.BindJSON(&category); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	err := t.todoSrv.AddCategory(ctx, category)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	ctx.JSON(http.StatusOK, " Category Added Successfully")
}

//EditTodo controller edit a todo
func (t todoCtrl) EditTodo(ctx *gin.Context) {
	var todo model.EditTodo
	if err := ctx.BindJSON(&todo); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	err := t.todoSrv.EditTodo(ctx, todo)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	ctx.JSON(http.StatusOK, " Todo updated Successfully")
}

//GetTodoByCategory to get todo by category
func (t todoCtrl) GetTodoByCategory(ctx *gin.Context) {
	category := ctx.Query("id")
	number, errParam := strconv.ParseUint(category, 10, 32)
	if errParam != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Please provide valid id")
		return
	}
	categoryInt := int(number)
	response, err := t.todoSrv.GetTodoByCategory(ctx, categoryInt)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	ctx.JSON(http.StatusOK, response)
}

//GetCategory controller to get category
func (t todoCtrl) GetCategory(ctx *gin.Context) {
	response, err := t.todoSrv.GetCategory(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	ctx.JSON(http.StatusOK, response)
}

//DeleteCategory controller to delete a category
func (t todoCtrl) DeleteCategory(ctx *gin.Context) {
	category := ctx.Query("id")
	number, errParam := strconv.ParseUint(category, 10, 32)
	if errParam != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "provide valid id")
		return
	}
	categoryInt := int(number)
	err := t.todoSrv.DeleteCategory(ctx, &categoryInt)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	ctx.JSON(http.StatusOK, "Successfully Deleted Category")
}
