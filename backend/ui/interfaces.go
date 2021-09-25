package ui

import "todo-tut/backend/entities"

type Service interface {
	GetAllTodos() ([]entities.Todo, error)
}