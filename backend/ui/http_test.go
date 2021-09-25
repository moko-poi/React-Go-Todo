package ui_test

import (
	"encoding/json"
	"fmt"
	"github.com/gomagedon/expectate"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-tut/backend/entities"
	"todo-tut/backend/ui"
)

// MockService

type MockService struct {
	err   error
	todos []entities.Todo
}

func (s MockService) GetAllTodos() ([]entities.Todo, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.todos, nil
}

// dummyTodos

var dummyTodos = []entities.Todo{
	{
		Title:       "todo 1",
		Description: "description of todo 1",
		IsCompleted: true,
	},
	{
		Title:       "todo 2",
		Description: "description of todo 2",
		IsCompleted: false,
	},
	{
		Title:       "todo 3",
		Description: "description of todo 3",
		IsCompleted: true,
	},
}

// HTTPTest

type HTTPTest struct {
	name string

	service *MockService
	inputMethod string
	inputURL string

	expectedStatus int
	expectedTodos []entities.Todo
}

func TestHTTP(t *testing.T) {

	tests := getTests()
	tests = append(tests, getDisallowedMethodTests()...)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testHTTP(t, test)
		})
	}
}

// Main Test

func testHTTP(t *testing.T, test HTTPTest) {
	expect := expectate.Expect(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(test.inputMethod, test.inputURL, nil)

	server := ui.NewHTTP()
	server.UseService(test.service)

	server.ServerHTTP(w, r)

	var body []entities.Todo
	json.NewDecoder(w.Result().Body).Decode(&body)

	expect(w.Result().StatusCode).ToBe(test.expectedStatus)
	expect(body).ToEqual(test.expectedTodos)
}

func getTests() []HTTPTest {
	return []HTTPTest{
		{
			name: "Random error gives 500 status and no tpdos",

			service:     &MockService{err: fmt.Errorf("something bad happened")},
			inputMethod: "GET",
			inputURL:    "http://mywebsite.com/todos/",

			expectedStatus: 500,
		},
		{
			name: "Random error gives 500 status and no tpdos",

			service:     &MockService{err: fmt.Errorf("something bad happened")},
			inputMethod: "GET",
			inputURL:    "http://mywebsite.com/todos",

			expectedStatus: 500,
		},
		{
			name: "Wrong path gives 404 status and no todos",

			service:     &MockService{todos: dummyTodos},
			inputMethod: "GET",
			inputURL:    "http://mywebsite.com/foo",

			expectedStatus: 404,
		},
		{
			name: "Wrong path gives 404 status and no todos",

			service:     &MockService{todos: dummyTodos},
			inputMethod: "GET",
			inputURL:    "http://mywebsite.com/bar",

			expectedStatus: 404,
		},
		{
			name: "Returns todos from service if no error",

			service: &MockService{todos: dummyTodos},
			inputMethod: "GET",
			inputURL: "http://mywebsite.com/todos/",

			expectedStatus: 200,
			expectedTodos: dummyTodos,
		},
	}
}

func getDisallowedMethodTests() []HTTPTest {
	tests := []HTTPTest{}

	disallowedMethods := []string {
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	}
	for _, method := range disallowedMethods {
		tests = append(tests, HTTPTest{
			name: fmt.Sprintf("Method %s gives 405 status and no todos", method),

			service: &MockService{todos: dummyTodos},
			inputURL: "http://mywebsite,com/todos/",
			inputMethod: method,

			expectedStatus: http.StatusMethodNotAllowed,
		})
	}

	return tests
}