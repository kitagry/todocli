package ui

import (
	"github.com/gdamore/tcell"
	"github.com/kitagry/go-todotxt"
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

	todolist []*todotxt.Task
}

// NewApplication returns UI.
func NewApplication(todolist []*todotxt.Task) *App {
	p := tview.NewPages()

	t := NewTable()
	t.WriteTasks(todolist)
	t.Select(1, 1).SetSelectable(true, false)
	t.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		inputText := tview.NewInputField().SetText(todolist[row-1].Description())
		inputText.SetDoneFunc(func(key tcell.Key) {
			todolist[row-1].SetDescription(inputText.GetText())
			t.WriteTask(todolist[row-1], row)
			p.RemovePage("input")
		})
		p.AddAndSwitchToPage("input", inputText, true)
	})

	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModNone {
			switch event.Rune() {
			case 'a':
				row, _ := t.GetSelection()
				todo := todolist[row-1]
				todo.SetPriority('A')
				t.WriteTask(todo, row)
			case 'b':
				row, _ := t.GetSelection()
				todo := todolist[row-1]
				todo.SetPriority('B')
				t.WriteTask(todo, row)
			case 'c':
				row, _ := t.GetSelection()
				todo := todolist[row-1]
				todo.SetPriority('C')
				t.WriteTask(todo, row)
			case 'd':
				row, _ := t.GetSelection()
				if row == 0 {
					return event
				}
				removeTask(todolist, row-1)
				t.RemoveRow(row)
			case 'x':
				row, _ := t.GetSelection()
				todo := todolist[row-1]
				if !todo.Completed {
					todo.Complete()
				} else {
					todo.Reopen()
				}
				t.WriteTask(todo, row)
			case 'n':
				inputText := tview.NewInputField()
				inputText.SetDoneFunc(func(key tcell.Key) {
					todo := todotxt.NewTask()
					todo.SetDescription(inputText.GetText())
					todolist = append(todolist, todo)
					t.WriteTask(todo, len(todolist))
					p.RemovePage("input")
				})
				p.AddAndSwitchToPage("input", inputText, true)
			}
		}
		return event
	})

	app := &App{
		Application: tview.NewApplication(),
		Pages:       p,
		Table:       t,
		todolist:    todolist,
	}

	p.AddPage("table", t, true, true)
	app.SetRoot(p, true).SetFocus(p)
	return app
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
