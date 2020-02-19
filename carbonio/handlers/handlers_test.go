package handlers

import "testing"

func TestName(t *testing.T) {
	hs := map[string]bool{}
	for _, h := range hndlrs {
		name := h.Name()
		if _, ok := hs[name]; ok {
			t.Errorf("duplicate handler name %q", name)
			continue
		}
		hs[name] = true
	}
}
