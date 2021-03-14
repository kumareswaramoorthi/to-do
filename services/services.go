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

func NewTodoService(todoDb database.TodoDatabase) TodoService {
	return todoService{
		todoDatabase: todoDb,
	}
}

// SignUp method creates a new user, it saves  user data in db
func (ds todoService) SignUp(ctxt *gin.Context, user *model.User) error {
	user.Prepare()
	// check the email format is valid or not
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return err
	}
	//check if email already exists
	EmailExists := ds.todoDatabase.CheckEmailExists(user.Email)
	if EmailExists {
		return errors.New("user already exists")
	}
	//Hash the password before saving
	user.HashBeforeSave()
	errCreate := ds.todoDatabase.CreateUser(user)
	if errCreate != nil {
		return errCreate
	}
	return nil
}

//SignIn methods gives the jwt token and user details provided if user credentials are valid.
func (ds todoService) SignIn(ctxt *gin.Context, user *model.UserLogin) (*model.SignInResponse, error) {
	//fetch the user details from database
	getuser, err := ds.todoDatabase.FindUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	// verify the password from database with incoming passwors in request
	err = utils.VerifyPassword(getuser.Password, user.Password)
	if err != nil {
		return nil, err
	}
	// if credentials are validated, create jwt token for the user
	token, err := auth.CreateToken(getuser.ID)
	if err != nil {
		return nil, err
	}
	// return the response
	return &model.SignInResponse{
		Token: token,
		User:  *getuser,
	}, nil
}

//AddTodo method is used to add new todo
func (ds todoService) AddTodo(ctxt *gin.Context, todo *model.Todo) error {
	// fetch the user id from gin context
	id, _ := ctxt.Get("user-id")
	todo.UserId = id.(int)
	// check the user has the category which is mentioned in the request data
	if todo.Category != 0 {
		category_user_id, err := ds.todoDatabase.GetCategoryUserById(todo.Category)
		if err != nil {
			return errors.New("category not found for this user")
		}
		//if the category is not valid, return error
		if id != *category_user_id {
			return errors.New("invalid category ID")
		}
	}
	// if the request data are all valid, add the todo and save it in the database
	err := ds.todoDatabase.AddTodo(todo)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

//DeleteTodo method deletes a todo
func (ds todoService) DeleteTodo(ctxt *gin.Context, todo string) error {
	// fetch the user-id from context
	id, _ := ctxt.Get("user-id")
	// check the todo id
	if todo == "" {
		return errors.New("todo_id nil")
	}
	//fetch the todo from database
	getTodo, err := ds.todoDatabase.GetTodoById(todo)
	if err != nil {
		return errors.New("todo does not exist")
	}
	if getTodo == nil {
		return errors.New("todo not found")
	}
	//check if the todo belongs to the current user
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

//GetAlltodo method fetches all todo that belongs to the current user
func (ds todoService) GetAlltodo(ctxt *gin.Context) (*[]model.Todo, error) {
	// fetch the user-id from context
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

//MarkTodo method used to mark the status of the todo
func (ds todoService) MarkTodo(ctxt *gin.Context, todoToMark model.MarkTodo) error {
	// fetch the user-id from context
	id, _ := ctxt.Get("user-id")
	var todo *model.Todo
	var completed int
	// get the dod from database
	todo, err := ds.todoDatabase.GetTodoById(fmt.Sprint(todoToMark.ID))
	if err != nil {
		return errors.New("unable to identify todo")
	}
	if todo == nil {
		return errors.New("todo does not exist")
	}
	// check if the todo belongs to the current user
	if todo.UserId != id {
		return errors.New("not Authorized to mark this todo")
	}
	// convert bool to int
	if todoToMark.Completed {
		completed = 1
	} else {
		completed = 0
	}
	// update status in database
	err = ds.todoDatabase.UpdateCompleted(completed, todoToMark.ID)
	if err != nil {
		return errors.New("unable to mark todo")
	}
	return nil

}

//AddCategory method used to add a category
func (ds todoService) AddCategory(ctxt *gin.Context, category model.Category) error {
	// fetch the user-id from context
	id, _ := ctxt.Get("user-id")
	category.UserId = id.(int)
	err := ds.todoDatabase.AddCategory(&category)
	if err != nil {
		return errors.New("unable to add category")
	}
	return nil
}

//EditTodo method used to edit a todo
func (ds todoService) EditTodo(ctxt *gin.Context, todoInput model.EditTodo) error {
	// check for the todo id present in request
	if todoInput.ID == 0 {
		return errors.New("please provide todo id")
	}
	// fetch the user id from gin context
	id, _ := ctxt.Get("user-id")
	// fetch the todo from database
	todo, err := ds.todoDatabase.GetTodoById(fmt.Sprint(todoInput.ID))
	if err != nil {
		return errors.New("unable to process todo")
	}
	// check if the todo belongs to the current user
	if todo.UserId != id {
		return errors.New("not Authorized to edit this todo")
	}
	//map the remaining fields from the database with todo from request
	editTodoPayload := utils.EditTodoMap(todoInput, *todo)
	// update database
	err = ds.todoDatabase.UpdateTodo(&editTodoPayload)
	if err != nil {
		return errors.New("unable to edit todo")
	}
	return nil
}

//GetTodoByCategory fetches the todos based on the category provided in request
func (ds todoService) GetTodoByCategory(ctxt *gin.Context, category_id int) (*[]model.Todo, error) {
	//fetch the  user-id from request
	id, _ := ctxt.Get("user-id")
	// fetch the user-id of that particular category to check if it belongs to current user
	category_user_id, err := ds.todoDatabase.GetCategoryUserById(int(category_id))
	if err != nil {
		return nil, errors.New("todo not found")
	}
	// if the category mentioned in request does not belong to the current user
	if id != *category_user_id {
		return nil, errors.New("not Authorized to view this todos")
	}
	//fetch todos based on the same type of category
	todos, err := ds.todoDatabase.GetAllTodoByCategory(id.(int), category_id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	if len(*todos) == 0 {
		return nil, errors.New(" no todos found for this category ")
	}
	return todos, nil
}

//GetCategory fetches all the category that belongs to the logged in user
func (ds todoService) GetCategory(ctxt *gin.Context) (*[]model.Category, error) {
	//fetch the user-id from context
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

//DeleteCategory method deletes the category if provided a valid category id
func (ds todoService) DeleteCategory(ctxt *gin.Context, category *int) error {
	// get the user id from context
	id, _ := ctxt.Get("user-id")
	// check the category id is present or not
	if category == nil {
		return errors.New("category id nil")
	}
	// get the category and store it in a variable
	getCategory, err := ds.todoDatabase.GetCategoryById(id.(int), *category)
	if err != nil {
		return errors.New("category does not exist for this user")
	}
	if getCategory == nil {
		return errors.New(" category not found")
	}
	//check if the category actually belongs to the curent user
	if id == getCategory.UserId {
		effect, err := ds.todoDatabase.DeleteCategory(*category)
		if err != nil {
			return errors.New("unable to delete category")
		}
		//if not deleted
		if effect == 0 {
			return errors.New("unable to delete category")
		}

	} else {
		return errors.New("not Authorized to delete this category")
	}
	return nil
}
