package usecase_test

import (
	"fmt"
	"github.com/gomagedon/expectate"
	"testing"
	"todo-tut/backend/entities"
	"todo-tut/backend/usecase"
)

var dummyTodos = []entities.Todo{
	{
		Title: "todo 1",
		Description: "description of todo 1",
		IsCompleted: true,
	},
	{
		Title: "todo 2",
		Description: "description of todo 2",
		IsCompleted: false,
	},
	{
		Title: "todo 3",
		Description: "description of todo 3",
		IsCompleted: true,
	},
}

type BadTodosRepo struct {}

func (BadTodosRepo) GetAllTodos() ([]entities.Todo, error) {
	return nil, fmt.Errorf("someshing went wrong")
}

type MockTodosRepo struct {}

func (MockTodosRepo) GetAllTodos() ([]entities.Todo, error) {
	return dummyTodos, nil
}

func TestGetTodos(t *testing.T) {
	// Test
	t.Run("Returns ErrInternal when TodosRepository returns error", func(t *testing.T) {
		expect := expectate.Expect(t)

		repo := new(BadTodosRepo)

		todos, err := usecase.GetTodos(repo)

		expect(err).ToBe(usecase.ErrInternal)
		if todos != nil {
			t.Fatalf("Expected todos to be nil")
		}
	})

	// Test
	t.Run("Returns todos from TodosRepository", func(t *testing.T) {
		expect := expectate.Expect(t)

		repo := new(MockTodosRepo)

		todos, err := usecase.GetTodos(repo)

		expect(err).ToBe(nil)
		expect(todos).ToEqual(dummyTodos)
	})
}
