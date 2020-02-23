package handlers

import "testing"

func TestHandlers(t *testing.T) {
	if _, err := NewHandlers(nil); err == nil {
		t.Errorf("expected an error")
	}
}
