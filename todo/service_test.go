package todo

import (
	"fmt"
	"testing"

	"github.com/kitagry/go-todotxt"
)

const taskLength = 3

func newTask(description string) *todotxt.Task {
	task := todotxt.NewTask()
	task.SetDescription(description)
	return task
}

func setupService() *Service {
	todolist := make([]*todotxt.Task, taskLength)
	for i := 0; i < taskLength; i++ {
		todolist[i] = newTask(fmt.Sprintf("task%d", i+1))
	}
	s := Service{
		todolist: todolist,
	}
	return &s
}

func TestService_SetPriority(t *testing.T) {
	s := setupService()
	tests := []struct {
		index    int
		priority byte
		haserror bool
	}{
		{
			index:    -1,
			priority: 'A',
			haserror: true,
		},
		{
			index:    0,
			priority: 'A',
			haserror: false,
		},
		{
			index:    taskLength - 1,
			priority: 'A',
			haserror: false,
		},
		{
			index:    taskLength,
			priority: 'A',
			haserror: true,
		},
	}

	for _, test := range tests {
		_, err := s.SetPriority(test.priority, test.index)
		if test.haserror {
			if err == nil {
				t.Errorf("s.SetPriority should return error")
			}
		} else {
			if err != nil {
				t.Errorf("s.SetPriority shouldn't return error: %v", err)
				continue
			}
			p := s.todolist[test.index].Priority()
			if p != test.priority {
				t.Errorf("Priority expect %b\ngot %b", test.priority, p)
			}
		}
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		index    int
		length   int
		haserror bool
	}{
		{
			index:    -1,
			length:   taskLength,
			haserror: true,
		},
		{
			index:    0,
			length:   taskLength - 1,
			haserror: false,
		},
		{
			index:    taskLength - 1,
			length:   taskLength - 1,
			haserror: false,
		},
		{
			index:    taskLength,
			length:   taskLength,
			haserror: true,
		},
	}

	for _, test := range tests {
		s := setupService()
		err := s.Delete(test.index)
		if (err != nil) != test.haserror {
			t.Errorf("Delete should return error: %v", err)
		}

		if len(s.todolist) != test.length {
			t.Errorf("Delete failed, len(todolist) expect %d\ngot %d", test.length, len(s.todolist))
		}
	}
}
