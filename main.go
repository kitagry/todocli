package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/kitagry/go-todotxt"
	"github.com/rivo/tview"
)

var todolist todotxt.TaskList

const headerLine = 0

func getColor(todo *todotxt.Task) tcell.Color {
	switch {
	case todo.Completed:
		return tcell.ColorGray
	case todo.IsOverdue():
		return tcell.ColorRed
	case todo.Priority == "A":
		return tcell.ColorBlue
	case todo.Priority == "B":
		return tcell.ColorAqua
	case todo.Priority == "C":
		return tcell.ColorPurple
	}
	return tcell.ColorWhite
}

func writeToTable(table *tview.Table, index int, row int) {
	todo := todolist[index]
	tableColor := getColor(todo)

	table.SetCell(row, 1, tview.NewTableCell(todo.Priority).SetTextColor(tableColor))
	table.SetCell(row, 2, tview.NewTableCell(todo.Todo).SetTextColor(tableColor))

	createdAt := ""
	if todo.HasCreatedDate() {
		createdAt = todo.CreatedDate.Format("2006-01-02")
	}
	table.SetCell(row, 3, tview.NewTableCell(createdAt).SetTextColor(tableColor))

	completedAt := ""
	if todo.HasCompletedDate() {
		completedAt = todo.CompletedDate.Format("2006-01-02")
	}
	table.SetCell(row, 4, tview.NewTableCell(completedAt).SetTextColor(tableColor))
}

func newTable(pages *tview.Pages) {
	table := tview.NewTable()
	// Set Header
	table.SetCell(headerLine, 1, tview.NewTableCell("Priority").SetTextColor(tcell.ColorYellow))
	table.SetCell(headerLine, 2, tview.NewTableCell("Task").SetTextColor(tcell.ColorYellow))
	table.SetCell(headerLine, 3, tview.NewTableCell("CreatedAt").SetTextColor(tcell.ColorYellow))
	table.SetCell(headerLine, 4, tview.NewTableCell("CompletedAt").SetTextColor(tcell.ColorYellow))

	for index := 0; index < len(todolist); index++ {
		writeToTable(table, index, index+1)
	}
	table.Select(1, 1).SetSelectable(true, false)

	table.SetSelectedFunc(func(row int, column int) {
		if row == 0 {
			return
		}
		inputText := tview.NewInputField().SetText(todolist[row-1].Todo)
		inputText.SetDoneFunc(func(key tcell.Key) {
			todolist[row-1].Todo = inputText.GetText()
			writeToTable(table, row-1, row)
			pages.RemovePage("input")
		})
		pages.AddAndSwitchToPage("input", inputText, true)
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModNone {
			switch event.Rune() {
			case 'a':
				row, _ := table.GetSelection()
				todo := todolist[row-1]
				todo.Priority = "A"
				writeToTable(table, row-1, row)
			case 'b':
				row, _ := table.GetSelection()
				todo := todolist[row-1]
				todo.Priority = "B"
				writeToTable(table, row-1, row)
			case 'c':
				row, _ := table.GetSelection()
				todo := todolist[row-1]
				todo.Priority = "C"
				writeToTable(table, row-1, row)
			case 'd':
				row, _ := table.GetSelection()
				todolist.RemoveTask(todolist[row-1])
				table.RemoveRow(row)
			case 'x':
				row, _ := table.GetSelection()
				todo := todolist[row-1]
				if todo.Completed {
					todo.Reopen()
				} else {
					todo.Complete()
				}
				writeToTable(table, row-1, row)
			case 'n':
				inputText := tview.NewInputField()
				inputText.SetDoneFunc(func(key tcell.Key) {
					todo := todotxt.NewTask(inputText.GetText())
					todo.Priority = "C"
					todolist.AddTask(todo)
					writeToTable(table, len(todolist)-1, len(todolist))
					pages.RemovePage("input")
				})
				pages.AddAndSwitchToPage("input", inputText, true)
			}
		}
		return event
	})

	pages.AddPage("todo-list", table, true, true)
}

func main() {
	f, err := os.Open("todo.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	todolist, err = todotxt.LoadFromFile(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	app := tview.NewApplication()

	pages := tview.NewPages()
	newTable(pages)

	app.SetRoot(pages, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if page, _ := pages.GetFrontPage(); page == "todo-list" && event.Modifiers() == tcell.ModNone {
			switch event.Rune() {
			case 'q':
				f, err := os.Create("todo.txt")
				if err != nil {
					return event
				}
				defer f.Close()

				todolist.WriteToFile(f)
				app.Stop()
			}
		}
		return event
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
