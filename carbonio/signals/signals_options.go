package signals

import "fmt"

type options struct {
	num    uint // Signal number (1-based, i.e. 1 is the first signal number).
	maxNum uint // Maximum number of signals of this type.

	dir  Dir
	conn Conn
	fmt  Fmt
	lvl  Lvl
}

func (o *options) validate() error {
	if o.maxNum == 0 {
		return fmt.Errorf("MaxNumber option unset")
	}
	if o.num <= 0 || o.num > o.maxNum {
		return fmt.Errorf("Number option out of range (0-%d)", o.maxNum)
	}
	return nil
}

func Number(v uint) func(*options) error {
	return func(o *options) error { return o.setNumber(v) }
}
func (o *options) setNumber(v uint) error {
	o.num = v
	return nil
}

func MaxNumber(v uint) func(*options) error {
	return func(o *options) error { return o.setMaxNumber(v) }
}
func (o *options) setMaxNumber(v uint) error {
	o.maxNum = v
	return nil
}

func Connector(v Conn) func(*options) error {
	return func(o *options) error { return o.setConnector(v) }
}
func (o *options) setConnector(v Conn) error {
	o.conn = v
	return nil
}

func Direction(v Dir) func(*options) error {
	return func(o *options) error { return o.setDirection(v) }
}
func (o *options) setDirection(v Dir) error {
	o.dir = v
	return nil
}

func Format(v Fmt) func(*options) error {
	return func(o *options) error { return o.setFormat(v) }
}
func (o *options) setFormat(v Fmt) error {
	o.fmt = v
	return nil
}

func Level(v Lvl) func(*options) error {
	return func(o *options) error { return o.setLevel(v) }
}
func (o *options) setLevel(v Lvl) error {
	o.lvl = v
	return nil
}
