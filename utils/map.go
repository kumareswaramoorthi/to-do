package utils

import "todo/model"

func EditTodoMap(todoInput model.EditTodo, todo model.Todo) model.EditTodo {
	if todoInput.Title == "" {
		todoInput.Title = todo.Title
	}
	if todoInput.Description == "" {
		todoInput.Description = todo.Description
	}
	if todoInput.DueDate == "" {
		todoInput.DueDate = todo.DueDate
	}
	if todoInput.Priority == "" {
		todoInput.Priority = todo.Priority
	}
	if todoInput.Completed == nil {
		todoInput.Completed = &todo.Completed
	}
	if todoInput.Category == 0 {
		todoInput.Category = todo.Category
	}
	todoInput.UserId = todo.UserId
	return todoInput
}
