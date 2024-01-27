package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Todo struct {
	Text      string
	Completed bool
}

// TodoApp struct represents the main application.
type TodoApp struct {
	app.Compo

	todos     []Todo
	newTodo   string
	filter    string
	allActive bool
}

// Render renders the TodoApp component.
func (a *TodoApp) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Todo App"),
		app.Input().
			Type("text").
			Placeholder("What needs to be done?").
			Value(a.newTodo).
			OnChange(a.onChangeNewTodo),
		app.Button().Text("Add").OnClick(a.addTodo),
		app.Br(),
		app.Select().
			OnChange(a.onChangeFilter).
			Body(
				app.Option().Value("all").Text("All"),
				app.Option().Value("active").Text("Active"),
				app.Option().Value("completed").Text("Completed"),
			),
		app.Button().Text("Clear Completed").OnClick(a.clearCompleted),
		app.Br(),
		app.Div().Body(
			app.Range(a.filteredTodos()).Slice(func(i int) app.UI {
				todo := a.todos[i]
				return app.Div().Body(
					app.Input().
						Type("checkbox").
						Checked(todo.Completed).
						OnChange(func(ctx app.Context, e app.Event) {
							a.toggleCompleted(i)
						}),
					app.Span().Text(todo.Text),
					app.Button().Text("Remove").OnClick(func(ctx app.Context, e app.Event) {
						a.removeTodo(i)
					}),
					app.Br(),
				)
			}),
		),
	)
}

func (a *TodoApp) onChangeNewTodo(ctx app.Context, e app.Event) {
	a.newTodo = ctx.JSSrc().Get("value").JSValue().String()
}

// addTodo adds a new todo to the list.
func (a *TodoApp) addTodo(ctx app.Context, e app.Event) {
	if a.newTodo != "" {
		a.todos = append(a.todos, Todo{Text: a.newTodo})
		a.newTodo = ""
	}
}

// onChangeFilter is called when the filter select field changes.
func (a *TodoApp) onChangeFilter(ctx app.Context, e app.Event) {
	a.filter = ctx.JSSrc().Get("value").JSValue().String()
}

// filteredTodos returns the todos based on the selected filter.
func (a *TodoApp) filteredTodos() []Todo {
	var filtered []Todo

	for _, todo := range a.todos {
		switch a.filter {
		case "all":
			filtered = append(filtered, todo)
		case "active":
			if !todo.Completed {
				filtered = append(filtered, todo)
			}
		case "completed":
			if todo.Completed {
				filtered = append(filtered, todo)
			}
		}
	}

	return filtered
}

// toggleCompleted toggles the completed status of a todo.
func (a *TodoApp) toggleCompleted(index int) {
	a.todos[index].Completed = !a.todos[index].Completed
}

// removeTodo removes a todo from the list.
func (a *TodoApp) removeTodo(index int) {
	a.todos = append(a.todos[:index], a.todos[index+1:]...)
}

// clearCompleted removes all completed todos from the list.
func (a *TodoApp) clearCompleted(ctx app.Context, e app.Event) {
	var active []Todo
	for _, todo := range a.todos {
		if !todo.Completed {
			active = append(active, todo)
		}
	}
	a.todos = active
}

func main() {
	todoApp := &TodoApp{
		todos:   []Todo{},
		newTodo: "",
		filter:  "all",
	}

	app.Route("/", todoApp)

	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
