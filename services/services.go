package services

import (
	"errors"
	"fmt"
	"todo/database"

	"todo/auth"
	"todo/model"
	"todo/utils"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
)

type TodoService interface {
	SignUp(ctxt *gin.Context, user *model.User) error
	SignIn(ctxt *gin.Context, user *model.UserLogin) (*model.SignInResponse, error)
	AddTodo(ctxt *gin.Context, todo *model.Todo) error
	DeleteTodo(ctxt *gin.Context, todo string) error
	GetAlltodo(ctxt *gin.Context) (*[]model.Todo, error)
	MarkTodo(ctxt *gin.Context, todoToMark model.MarkTodo) error
	AddCategory(ctxt *gin.Context, category model.Category) error
	EditTodo(ctxt *gin.Context, todoInput model.EditTodo) error
	GetTodoByCategory(ctxt *gin.Context, category_id int) (*[]model.Todo, error)
	GetCategory(ctxt *gin.Context) (*[]model.Category, error)
	DeleteCategory(ctxt *gin.Context, category *int) error
}

type todoService struct {
	todoDatabase database.TodoDatabase
}

func NewTodoService(todoDb database.TodoDatabase) todoService {
	return todoService{
		todoDatabase: todoDb,
	}
}

func (ds todoService) SignUp(ctxt *gin.Context, user *model.User) error {
	user.Prepare()
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return err
	}
	//check if email already exists
	EmailExists := ds.todoDatabase.CheckEmailExists(user.Email)
	if EmailExists {
		return errors.New("user already exists")
	}
	user.HashBeforeSave()
	errCreate := ds.todoDatabase.CreateUser(user)
	if errCreate != nil {
		return errCreate
	}
	return nil
}

func (ds todoService) SignIn(ctxt *gin.Context, user *model.UserLogin) (*model.SignInResponse, error) {
	getuser, err := ds.todoDatabase.FindUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	err = utils.VerifyPassword(getuser.Password, user.Password)
	if err != nil {
		return nil, err
	}
	token, err := auth.CreateToken(getuser.ID)
	if err != nil {
		return nil, err
	}

	return &model.SignInResponse{
		Token: token,
		User:  *getuser,
	}, nil
}

func (ds todoService) AddTodo(ctxt *gin.Context, todo *model.Todo) error {
	id, _ := ctxt.Get("user-id")
	todo.UserId = id.(int)
	if todo.Category != 0 {
		category_user_id, err := ds.todoDatabase.GetCategoryUserById(todo.Category)
		if err != nil {
			return errors.New("category not found for this user")
		}
		if id != *category_user_id {
			return errors.New("invalid category ID")
		}
	}
	err := ds.todoDatabase.AddTodo(todo)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds todoService) DeleteTodo(ctxt *gin.Context, todo string) error {
	id, _ := ctxt.Get("user-id")
	if todo == "" {
		return errors.New("todo_id nil")
	}
	getTodo, err := ds.todoDatabase.GetTodoById(todo)
	if err != nil {
		return errors.New("todo does not exist")
	}
	if getTodo == nil {
		return errors.New("todo not found")
	}
	if id == getTodo.UserId {
		effect, err := ds.todoDatabase.DeleteTodo(todo)
		if err != nil {
			return errors.New("unable to delete todo")
		}
		if effect == 0 {
			return errors.New("unable to delete todo")
		}

	} else {
		return errors.New("not Authorized to delete this todo")
	}
	return nil
}

func (ds todoService) GetAlltodo(ctxt *gin.Context) (*[]model.Todo, error) {
	id, _ := ctxt.Get("user-id")
	todos, err := ds.todoDatabase.GetAllTodo(id.(int))
	if err != nil {
		return nil, errors.New("unable to fetch todos")
	}
	if len(*todos) == 0 {
		return nil, errors.New("no todos found for this user")

	}
	return todos, nil
}

func (ds todoService) MarkTodo(ctxt *gin.Context, todoToMark model.MarkTodo) error {
	id, _ := ctxt.Get("user-id")
	var todo *model.Todo
	var completed int
	todo, err := ds.todoDatabase.GetTodoById(fmt.Sprint(todoToMark.ID))
	if err != nil {
		return errors.New("unable to identify todo")
	}
	if todo == nil {
		return errors.New("todo does not exist")
	}
	if todo.UserId != id {
		return errors.New("not Authorized to mark this todo")
	}
	if todoToMark.Completed {
		completed = 1
	} else {
		completed = 0
	}
	err = ds.todoDatabase.UpdateCompleted(completed, todoToMark.ID)
	if err != nil {
		return errors.New("unable to mark todo")
	}
	return nil

}

func (ds todoService) AddCategory(ctxt *gin.Context, category model.Category) error {
	id, _ := ctxt.Get("user-id")
	category.UserId = id.(int)
	err := ds.todoDatabase.AddCategory(&category)
	if err != nil {
		return errors.New("unable to add category")
	}
	return nil
}

func (ds todoService) EditTodo(ctxt *gin.Context, todoInput model.EditTodo) error {
	if todoInput.ID == 0 {
		return errors.New("please provide todo id")
	}
	id, _ := ctxt.Get("user-id")
	todo, err := ds.todoDatabase.GetTodoById(fmt.Sprint(todoInput.ID))
	if err != nil {
		return errors.New("unable to process todo")
	}
	if todo.UserId != id {
		return errors.New("not Authorized to edit this todo")
	}
	editTodoPayload := utils.EditTodoMap(todoInput, *todo)
	err = ds.todoDatabase.UpdateTodo(&editTodoPayload)
	if err != nil {
		return errors.New("unable to edit todo")
	}
	return nil
}

func (ds todoService) GetTodoByCategory(ctxt *gin.Context, category_id int) (*[]model.Todo, error) {
	id, _ := ctxt.Get("user-id")
	category_user_id, err := ds.todoDatabase.GetCategoryUserById(int(category_id))
	if err != nil {
		return nil, errors.New("todo not found")
	}
	if id != *category_user_id {
		return nil, errors.New("not Authorized to view this todos")
	}
	todos, err := ds.todoDatabase.GetAllTodoByCategory(id.(int), category_id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	if len(*todos) == 0 {
		return nil, errors.New(" no todos found for this category ")
	}
	return todos, nil
}

func (ds todoService) GetCategory(ctxt *gin.Context) (*[]model.Category, error) {
	id, _ := ctxt.Get("user-id")
	category, err := ds.todoDatabase.GetCategory(id.(int))
	if err != nil {
		return nil, errors.New("no category found for this user")
	}
	if len(*category) == 0 {
		return nil, errors.New(" no category found for this user")
	}
	return category, nil
}

func (ds todoService) DeleteCategory(ctxt *gin.Context, category *int) error {
	id, _ := ctxt.Get("user-id")
	if category == nil {
		return errors.New("category id nil")
	}
	getCategory, err := ds.todoDatabase.GetCategoryById(id.(int), *category)
	if err != nil {
		return errors.New("category does not exist ")
	}
	if getCategory == nil {
		return errors.New(" category not found")
	}
	if id == getCategory.UserId {
		effect, err := ds.todoDatabase.DeleteCategory(*category)
		if err != nil {
			return errors.New("unable to delete category")
		}
		if effect == 0 {
			return errors.New("unable to delete category")
		}

	} else {
		return errors.New("not Authorized to delete this category")
	}
	return nil
}
