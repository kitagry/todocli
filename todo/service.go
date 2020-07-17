package todo

import (
	"fmt"
	"os"
	"sort"

	"github.com/kitagry/go-todotxt"
	"golang.org/x/xerrors"
)

type Service struct {
	todolist []*todotxt.Task
}

func NewService(todolist []*todotxt.Task) *Service {
	return &Service{
		todolist: todolist,
	}
}

func (s *Service) AddNewTask(description string) *todotxt.Task {
	todo := todotxt.NewTask()
	todo.SetDescription(description)
	s.todolist = append(s.todolist, todo)
	return todo
}

func (s *Service) Length() int {
	return len(s.todolist)
}

func (s *Service) GetTask(index int) (*todotxt.Task, error) {
	if index < 0 || index >= len(s.todolist) {
		return nil, fmt.Errorf("index should from 0 to %d", len(s.todolist)-1)
	}
	return s.todolist[index], nil
}

func (s *Service) SetPriority(priority byte, index int) (*todotxt.Task, error) {
	if index < 0 || index >= len(s.todolist) {
		return nil, fmt.Errorf("index should from 0 to %d", len(s.todolist)-1)
	}
	todo := s.todolist[index]
	err := todo.SetPriority(priority)
	if err != nil {
		return nil, xerrors.Errorf("failed to SetPriority: %w", err)
	}
	return todo, nil
}

func (s *Service) Delete(index int) error {
	if index < 0 || index >= len(s.todolist) {
		return fmt.Errorf("index should from 0 to %d", len(s.todolist)-1)
	}
	s.todolist = removeTask(s.todolist, index)
	return nil
}

func (s *Service) ToggleCompleted(index int) (*todotxt.Task, error) {
	if index < 0 || index >= len(s.todolist) {
		return nil, fmt.Errorf("index should from 0 to %d", len(s.todolist)-1)
	}

	todo := s.todolist[index]
	if !todo.Completed {
		todo.Complete()
	} else {
		todo.Reopen()
	}
	return todo, nil
}

func (s *Service) SortPriorityAsc() {
	sort.SliceStable(s.todolist, func(i, j int) bool {
		if s.todolist[i].Completed {
			return false
		} else if s.todolist[j].Completed {
			return true
		}

		return s.todolist[i].Priority() > s.todolist[j].Priority()
	})
}

func (s *Service) SortPriorityDesc() {
	sort.SliceStable(s.todolist, func(i, j int) bool {
		if s.todolist[i].Completed {
			return false
		} else if s.todolist[j].Completed {
			return true
		}

		if s.todolist[i].Priority() == 0 {
			return false
		} else if s.todolist[j].Priority() == 0 {
			return true
		}
		return s.todolist[i].Priority() < s.todolist[j].Priority()
	})
}

func (s *Service) MoveCompletedTaskToBottom() {
	sort.SliceStable(s.todolist, func(i, j int) bool {
		return !s.todolist[i].Completed
	})
}

func (s *Service) SaveTodotxt(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return xerrors.Errorf("Failed to create %s: %w", filename, err)
	}
	defer f.Close()

	w := todotxt.NewWriter(f)
	return w.WriteAll(s.todolist)
}

func removeTask(list []*todotxt.Task, index int) []*todotxt.Task {
	if index == 0 {
		return list[1:]
	} else if index == len(list)-1 {
		return list[:index]
	}
	return append(list[:index], list[index+1:]...)
}
