package devices

import "testing"

func TestStage16_New(t *testing.T) {
	_, err := NewStage16()
	if err != nil {
		t.Fatalf("doh!")
	}
}
