package ui

import (
	"fmt"
	"unicode"

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

	t.SetInputCapture(app.EventHandler)

	p.AddPage("table", t, true, true)
	app.SetRoot(p, true).SetFocus(p)
	return app
}

func (a *App) SaveTodotxt(filename string) error {
	return a.service.SaveTodotxt(filename)
}

func (a *App) EventHandler(event *tcell.EventKey) *tcell.EventKey {
	t := a.Table
	p := a.Pages

	if event.Modifiers() == tcell.ModNone {
		switch event.Rune() {
		case 'a', 'b', 'c':
			row, _ := t.GetSelection()
			todo, err := a.service.SetPriority(byte(unicode.ToUpper(event.Rune())), row-1)
			if err == nil {
				t.WriteTask(todo, row)
			}
		case 'd':
			row, _ := t.GetSelection()
			todo, err := a.service.GetTask(row - 1)
			if err != nil {
				return event
			}
			confirm := tview.NewModal().
				SetText(fmt.Sprintf(`Do you want to delete task?\n"%s"`, todo.Description())).
				AddButtons([]string{"Delete", "Cancel"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Delete" {
						err = a.service.Delete(row - 1)
						if err == nil {
							t.RemoveRow(row)
						}
					}
					p.RemovePage("confirm")
				})
			p.AddAndSwitchToPage("confirm", confirm, true)
		case 's':
			a.AddSortListPage()
		case 'x':
			row, _ := t.GetSelection()
			todo, err := a.service.ToggleCompleted(row - 1)
			if err == nil {
				t.WriteTask(todo, row)
			}
		case 'n':
			inputText := tview.NewInputField().SetLabel("Input new task description: ")
			inputText.SetDoneFunc(func(key tcell.Key) {
				todo := a.service.AddNewTask(inputText.GetText())
				t.WriteTask(todo, a.service.Length())
				p.RemovePage("input")
			})
			p.AddAndSwitchToPage("input", inputText, true)
		}
	}
	return event
}

func (a *App) AddSortListPage() {
	list := tview.NewList().
		AddItem("Sort by priority desc", "A to Z", 'a', func() {
			a.service.SortPriorityDesc()
			a.Table.WriteTasks(a.service)
			a.Pages.RemovePage("sort")
		}).
		AddItem("Sort by priority asc", "Z to A", 'b', func() {
			a.service.SortPriorityAsc()
			a.Table.WriteTasks(a.service)
			a.Pages.RemovePage("sort")
		})
	a.Pages.AddAndSwitchToPage("sort", list, true)
}
