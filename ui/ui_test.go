package ui

import (
	"testing"

	"github.com/kitagry/go-todotxt"
)

func TestNewUi(t *testing.T) {
	u := NewApplication([]*todotxt.Task{})
	if u.Pages == nil {
		t.Errorf("u.Pages should be not null")
	}
}
