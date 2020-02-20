package ui

import "testing"

func TestNewUi(t *testing.T) {
	u := NewUI()
	if u.Pages == nil {
		t.Errorf("u.Pages should be not null")
	}
}
