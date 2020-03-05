package ui

import (
	"fmt"
	"os"
	"sort"

	"github.com/gdamore/tcell"
	"github.com/kitagry/go-todotxt"
	"github.com/rivo/tview"
	"golang.org/x/xerrors"
)

const (
	tableName = "todo-list"
)

// App is the whole view of this plugin.
type App struct {
	*tview.Application
	Pages *tview.Pages
	Table *Table

	todolist []*todotxt.Task
}

// NewApplication returns UI.
func NewApplication(todolist []*todotxt.Task) *App {
	p := tview.NewPages()

	t := NewTable()

	app := &App{
		Application: tview.NewApplication(),
		Pages:       p,
		Table:       t,
		todolist:    todolist,
	}

	t.WriteTasks(app.todolist)
	t.Select(1, 1).SetSelectable(true, false)
	t.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		inputText := tview.NewInputField().SetText(app.todolist[row-1].Description())
		inputText.SetDoneFunc(func(key tcell.Key) {
			app.todolist[row-1].SetDescription(inputText.GetText())
			t.WriteTask(app.todolist[row-1], row)
			p.RemovePage("input")
		})
		p.AddAndSwitchToPage("input", inputText, true)
	})

	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModNone {
			switch event.Rune() {
			case 'a':
				row, _ := t.GetSelection()
				todo := app.todolist[row-1]
				todo.SetPriority('A')
				t.WriteTask(todo, row)
			case 'b':
				row, _ := t.GetSelection()
				todo := app.todolist[row-1]
				todo.SetPriority('B')
				t.WriteTask(todo, row)
			case 'c':
				row, _ := t.GetSelection()
				todo := app.todolist[row-1]
				todo.SetPriority('C')
				t.WriteTask(todo, row)
			case 'd':
				row, _ := t.GetSelection()
				if row == 0 {
					return event
				}
				confirm := tview.NewModal().
					SetText(fmt.Sprintf(`Do you want to delete task?\n"%s"`, app.todolist[row-1].Description())).
					AddButtons([]string{"Delete", "Cancel"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "Delete" {
							app.todolist = removeTask(app.todolist, row-1)
							t.RemoveRow(row)
						}
						p.RemovePage("confirm")
					})
				p.AddAndSwitchToPage("confirm", confirm, true)
			case 's':
				list := tview.NewList().
					AddItem("Sort by priority desc", "A to Z", 'a', func() {
						sort.Slice(app.todolist, func(i, j int) bool {
							return app.todolist[i].Priority() < app.todolist[j].Priority()
						})
						t.WriteTasks(app.todolist)
						p.RemovePage("sort")
					}).
					AddItem("Sort by priority asc", "Z to A", 'b', func() {
						sort.Slice(app.todolist, func(i, j int) bool {
							return app.todolist[i].Priority() > app.todolist[j].Priority()
						})
						t.WriteTasks(app.todolist)
						p.RemovePage("sort")
					}).
					AddItem("Move done task to bottom", "", 'c', func() {
						sort.Slice(app.todolist, func(i, j int) bool {
							return !app.todolist[i].Completed
						})
						t.WriteTasks(app.todolist)
						p.RemovePage("sort")
					})
				p.AddAndSwitchToPage("sort", list, true)
			case 'x':
				row, _ := t.GetSelection()
				todo := app.todolist[row-1]
				if !todo.Completed {
					todo.Complete()
				} else {
					todo.Reopen()
				}
				t.WriteTask(todo, row)
			case 'n':
				inputText := tview.NewInputField().SetLabel("Input new task description: ")
				inputText.SetDoneFunc(func(key tcell.Key) {
					todo := todotxt.NewTask()
					todo.SetDescription(inputText.GetText())
					app.todolist = append(app.todolist, todo)
					t.WriteTask(todo, len(app.todolist))
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
	f, err := os.Create(filename)
	if err != nil {
		return xerrors.Errorf("Failed to create %s: %w", filename, err)
	}
	defer f.Close()

	w := todotxt.NewWriter(f)
	return w.WriteAll(a.todolist)
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
