// Code generated by "stringer -type=LEDState ledstates.go"; DO NOT EDIT.

package leds

import "strconv"

const (
	_LEDState_name_0 = "UnknownOffAlertOn"
	_LEDState_name_1 = "testLEDState"
)

var (
	_LEDState_index_0 = [...]uint8{0, 7, 10, 15, 17}
)

func (i LEDState) String() string {
	switch {
	case 0 <= i && i <= 3:
		return _LEDState_name_0[_LEDState_index_0[i]:_LEDState_index_0[i+1]]
	case i == 255:
		return _LEDState_name_1
	default:
		return "LEDState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
