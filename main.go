package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/kitagry/go-todotxt"
	"github.com/kitagry/todocli/ui"
)

var (
	todotxtPath string

	todolist []*todotxt.Task
)

const headerLine = 0

func main() {
	flag.StringVar(&todotxtPath, "f", "todo.txt", "todo.txt path")
	flag.Parse()

	f, err := os.Open(todotxtPath)
	if err == nil {
		defer f.Close()
		r := todotxt.NewReader(f)
		todolist, err = r.ReadAll()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	app := ui.NewApplication(todolist)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if page, _ := app.Pages.GetFrontPage(); page == "table" && event.Modifiers() == tcell.ModNone {
			switch event.Rune() {
			case 'q':
				err := app.SaveTodotxt(todotxtPath)
				if err != nil {
					fmt.Println(err)
				}
				app.Stop()
			case 'w':
				err := app.SaveTodotxt(todotxtPath)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		return event
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
