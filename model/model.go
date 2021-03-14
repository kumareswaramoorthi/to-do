package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"  binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"  binding:"required"`
	Priority    string `json:"priority"`
	Completed   bool   `json:"completed"`
	UserId      int    `json:"userId"`
	Category    int    `json:"category"`
}
type MarkTodo struct {
	ID        int  `json:"id" binding:"required"`
	Completed bool `json:"completed" `
}

type Category struct {
	ID     int    `json:"id"`
	Name   string `json:"name" binding:"required"`
	UserId int    `json:"userId"`
}

type Id struct {
	ID int `json:"id"`
}

type SignInResponse struct {
	Token string
	User  User
}
type EditTodo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	Priority    string `json:"priority"`
	Completed   *bool  `json:"completed"`
	UserId      int    `json:"userId"`
	Category    int    `json:"category"`
}

// Hash the password before saving into database
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Verify the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func (u *User) HashBeforeSave() error {
	hashed, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}
func (u *User) Prepare() {
	u.ID = 0
}
