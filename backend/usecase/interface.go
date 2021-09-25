package usecase

import "todo-tut/backend/entities"

type TodosRepository interface {
	GetAllTodos() ([]entities.Todo, error)
}