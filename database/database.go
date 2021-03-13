package database

import (
	"database/sql"
	"errors"
	"fmt"

	"todo/model"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlCreateUser = `
    CREATE TABLE IF NOT EXISTS user(
        user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name VARCHAR NOT NULL,
		email VARCHAR NOT NULL,
		password VARCHAR NOT NULL
    );
    `
	sqlCreateCategory = `
    CREATE TABLE IF NOT EXISTS category(
        category_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        category_name VARCHAR NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES user (user_id)
    );
    `
	sqlCreateTodo = `
    CREATE TABLE IF NOT EXISTS todo(
        todo_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        title VARCHAR NOT NULL,
		description VARCHAR,
		due_date DATE,
		priority VARCHAR,
		completed INTEGER DEFAULT 0,
		user_id INTEGER NOT NULL,
		category INTEGER,
		FOREIGN KEY (user_id)  REFERENCES user (user_id) ,
		FOREIGN KEY (category)  REFERENCES category (category_id)
    );
    `
	sqlInsertUser = `
	INSERT INTO user 
		(name,email,password) VALUES (?,?,?)
	`
	sqlFindUserByEmail = `
	SELECT * FROM user
		WHERE email = ?;
	`
	sqlInsertTodo = `
	INSERT INTO todo
		(title,description,due_date,priority,completed,user_id,category)
		VALUES (?,?,?,?,?,?,?);
		`
	sqlInsertCategory = `
	INSERT INTO category
		(category_name,user_id)
		VALUES (?,?);
		`
	sqlGetCategory = `
	SELECT * FROM category 
		WHERE user_id = ?;
		`
	sqlGetCategoryById = `
	SELECT * FROM category 
		WHERE user_id = ? 
		AND category_id = ?;
		`
	sqlDeleteTodo = `
	DELETE from todo 
		WHERE todo_id = ? 
    `
	sqlDeleteCategory = `
	DELETE from category 
		WHERE category_id = ?
    `
	sqlUpdateTodo = `
	UPDATE todo 
		SET title = ?,
		description = ?,
		due_date = ?,
		priority = ?,
		completed = ?,
		user_id = ?,
		category = ?
		WHERE todo_id = ?
	`
	sqlUpdateTodoCompleted = `
	UPDATE todo 
		SET completed = ?
	 	WHERE todo_id = ?;
	`
	sqlGetAllTodo = `
	SELECT * FROM todo 
		WHERE user_id = ?
	 `
	sqlGetTodoById = `
	SELECT * FROM todo 
		WHERE todo_id = ?
	`

	sqlCategoryById = `
	SELECT user_id FROM category 
		WHERE category_id = ?
	`

	sqlGetAllTodoByCategory = `
	SELECT * FROM todo 
		WHERE user_id = ? AND category = ?
	 `
)

type TodoDatabase interface {
	CreateUser(u *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	AddTodo(to *model.Todo) error
	GetTodoById(id string) (*model.Todo, error)
	DeleteTodo(id string) (int64, error)
	GetAllTodo(id int) (*[]model.Todo, error)
	UpdateTodo(getTodo *model.EditTodo) error
	UpdateCompleted(completed int, id int) error
	AddCategory(category *model.Category) error
	GetCategoryUserById(id int) (*int, error)
	GetAllTodoByCategory(id int, categoryId int) (*[]model.Todo, error)
	CheckEmailExists(email string) bool
	GetCategory(id int) (*[]model.Category, error)
	DeleteCategory(id int) (int64, error)
	GetCategoryById(userId, categoryId int) (*model.Category, error)
}
type todoDatabase struct {
	db *sql.DB
}

func NewTodoDatabase(db *sql.DB) TodoDatabase {
	return todoDatabase{db: db}
}

func InitDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "./sqliteDB/todo.db")
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("no database found")
	}
	err = migrate(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {

	_, err := db.Exec(sqlCreateUser)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlCreateCategory)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlCreateTodo)
	if err != nil {
		return err
	}
	return nil
}

func (t todoDatabase) CreateUser(u *model.User) error {
	_, err := t.db.Exec(sqlInsertUser, u.Name, u.Email, u.Password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (t todoDatabase) FindUserByEmail(email string) (*model.User, error) {
	getuser := model.User{}
	err := t.db.QueryRow(sqlFindUserByEmail, email).Scan(&getuser.ID, &getuser.Name, &getuser.Email, &getuser.Password)
	if err != nil {
		return nil, err
	}
	return &getuser, nil
}

func (t todoDatabase) AddTodo(to *model.Todo) error {
	_, err := t.db.Exec(sqlInsertTodo, to.Title, to.Description, to.DueDate, to.Priority, to.Completed, to.UserId, to.Category)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (t todoDatabase) GetTodoById(id string) (*model.Todo, error) {
	getTodo := model.Todo{}
	err := t.db.QueryRow(sqlGetTodoById, id).Scan(&getTodo.ID, &getTodo.Title, &getTodo.Description, &getTodo.DueDate, &getTodo.Priority, &getTodo.Completed, &getTodo.UserId, &getTodo.Category)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(&getTodo)
	return &getTodo, nil
}

func (t todoDatabase) DeleteTodo(id string) (int64, error) {
	res, err := t.db.Exec(sqlDeleteTodo, id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	n, err := res.RowsAffected()
	if n != 1 {
		return 0, err
	}
	return n, nil
}

func (t todoDatabase) GetAllTodo(id int) (*[]model.Todo, error) {
	var TodoList []model.Todo
	var getTodo model.Todo
	rows, err := t.db.Query(sqlGetAllTodo, id)
	if err != nil {

		return nil, err
	}
	for rows.Next() {
		rows.Scan(&getTodo.ID, &getTodo.Title, &getTodo.Description, &getTodo.DueDate, &getTodo.Priority, &getTodo.Completed, &getTodo.UserId, &getTodo.Category)
		TodoList = append(TodoList, getTodo)
	}
	return &TodoList, nil
}

func (t todoDatabase) UpdateTodo(getTodo *model.EditTodo) error {
	_, err := t.db.Exec(sqlUpdateTodo, &getTodo.Title, &getTodo.Description, &getTodo.DueDate, &getTodo.Priority, &getTodo.Completed, &getTodo.UserId, &getTodo.Category, &getTodo.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t todoDatabase) UpdateCompleted(completed int, id int) error {
	_, err := t.db.Exec(sqlUpdateTodoCompleted, &completed, &id)
	if err != nil {
		return err
	}
	return nil
}

func (t todoDatabase) AddCategory(category *model.Category) error {
	_, err := t.db.Exec(sqlInsertCategory, &category.Name, &category.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (t todoDatabase) GetCategoryUserById(id int) (*int, error) {
	var categoryUserId int
	err := t.db.QueryRow(sqlCategoryById, &id).Scan(&categoryUserId)
	if err != nil {
		return nil, err
	}
	return &categoryUserId, nil
}

func (t todoDatabase) GetAllTodoByCategory(id int, categoryId int) (*[]model.Todo, error) {
	var TodoList []model.Todo
	var getTodo model.Todo
	rows, err := t.db.Query(sqlGetAllTodoByCategory, id, categoryId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&getTodo.ID, &getTodo.Title, &getTodo.Description, &getTodo.DueDate, &getTodo.Priority, &getTodo.Completed, &getTodo.UserId, &getTodo.Category)
		TodoList = append(TodoList, getTodo)
	}
	return &TodoList, nil
}

func (t todoDatabase) CheckEmailExists(email string) bool {
	row := t.db.QueryRow("select email from user where email= ?", email)
	temp := ""
	row.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}

func (t todoDatabase) GetCategory(id int) (*[]model.Category, error) {
	var CategoryList []model.Category
	var category model.Category
	rows, err := t.db.Query(sqlGetCategory, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&category.ID, &category.Name, &category.UserId)
		CategoryList = append(CategoryList, category)
	}
	return &CategoryList, nil
}

func (t todoDatabase) DeleteCategory(id int) (int64, error) {
	res, err := t.db.Exec(sqlDeleteCategory, id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	n, err := res.RowsAffected()
	if n != 1 {
		return 0, err
	}
	return n, nil
}

func (t todoDatabase) GetCategoryById(userId, categoryId int) (*model.Category, error) {
	var category model.Category
	err := t.db.QueryRow(sqlGetCategoryById, userId, categoryId).Scan(&category.ID, &category.Name, &category.UserId)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
