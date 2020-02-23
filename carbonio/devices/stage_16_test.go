package devices

import "testing"

func TestStage16_New(t *testing.T) {
	_, err := NewStage16(
		SPIDelayRead(true),
	)
	if err != nil {
		t.Fatalf("error instantiating Stage16; %s", err)
	}
}
