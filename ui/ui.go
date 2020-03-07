package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/kitagry/go-todotxt"
	"github.com/kitagry/todocli/todo"
	"github.com/rivo/tview"
)

const (
	tableName = "todo-list"
)

// App is the whole view of this plugin.
type App struct {
	*tview.Application
	Pages *tview.Pages
	Table *Table

	service *todo.Service
}

// NewApplication returns UI.
func NewApplication(todolist []*todotxt.Task) *App {
	p := tview.NewPages()

	t := NewTable()

	app := &App{
		Application: tview.NewApplication(),
		Pages:       p,
		Table:       t,
		service:     todo.NewService(todolist),
	}

	t.WriteTasks(app.service)
	t.Select(1, 1).SetSelectable(true, false)
	t.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		todo, err := app.service.GetTask(row - 1)
		if err != nil {
			return
		}
		inputText := tview.NewInputField().SetText(todo.Description())
		inputText.SetDoneFunc(func(key tcell.Key) {
			todo.SetDescription(inputText.GetText())
			t.WriteTask(todo, row)
			p.RemovePage("input")
		})
		p.AddAndSwitchToPage("input", inputText, true)
	})

	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModNone {
			switch event.Rune() {
			case 'a':
				row, _ := t.GetSelection()
				todo, err := app.service.SetPriority('A', row-1)
				if err == nil {
					t.WriteTask(todo, row)
				}
			case 'b':
				row, _ := t.GetSelection()
				todo, err := app.service.SetPriority('B', row-1)
				if err == nil {
					t.WriteTask(todo, row)
				}
			case 'c':
				row, _ := t.GetSelection()
				todo, err := app.service.SetPriority('C', row-1)
				if err == nil {
					t.WriteTask(todo, row)
				}
			case 'd':
				row, _ := t.GetSelection()
				todo, err := app.service.GetTask(row - 1)
				if err != nil {
					return event
				}
				confirm := tview.NewModal().
					SetText(fmt.Sprintf(`Do you want to delete task?\n"%s"`, todo.Description())).
					AddButtons([]string{"Delete", "Cancel"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "Delete" {
							app.service.Delete(row - 1)
							t.RemoveRow(row)
						}
						p.RemovePage("confirm")
					})
				p.AddAndSwitchToPage("confirm", confirm, true)
			case 's':
				list := tview.NewList().
					AddItem("Sort by priority desc", "A to Z", 'a', func() {
						app.service.SortPriorityDesc()
						t.WriteTasks(app.service)
						p.RemovePage("sort")
					}).
					AddItem("Sort by priority asc", "Z to A", 'b', func() {
						app.service.SortPriorityAsc()
						t.WriteTasks(app.service)
						p.RemovePage("sort")
					}).
					AddItem("Move done task to bottom", "", 'c', func() {
						app.service.MoveCompletedTaskToBottom()
						t.WriteTasks(app.service)
						p.RemovePage("sort")
					})
				p.AddAndSwitchToPage("sort", list, true)
			case 'x':
				row, _ := t.GetSelection()
				todo, err := app.service.ToggleCompleted(row - 1)
				if err == nil {
					t.WriteTask(todo, row)
				}
			case 'n':
				inputText := tview.NewInputField().SetLabel("Input new task description: ")
				inputText.SetDoneFunc(func(key tcell.Key) {
					todo := app.service.AddNewTask(inputText.GetText())
					t.WriteTask(todo, app.service.Length())
					p.RemovePage("input")
				})
				p.AddAndSwitchToPage("input", inputText, true)
			}
		}
		return event
	})

	p.AddPage("table", t, true, true)
	app.SetRoot(p, true).SetFocus(p)
	return app
}

func (a *App) SaveTodotxt(filename string) error {
	return a.service.SaveTodotxt(filename)
}

func removeTask(list []*todotxt.Task, index int) []*todotxt.Task {
	if index < 0 || index >= len(list) {
		return list
	}

	if index == 0 {
		return list[1:]
	} else if index == len(list)-1 {
		return list[:index]
	}
	return append(list[:index], list[index+1:]...)
}
