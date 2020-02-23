package handlers

import "fmt"

type options struct {
	// Local flags.
	port int
	raw  bool
}

func (o *options) validate() error {
	if o.port == 0 {
		return fmt.Errorf("port option missing")
	}
	return nil
}

func Port(v int) func(*options) error {
	return func(o *options) error { return o.setPort(v) }
}
func (o *options) setPort(v int) error {
	o.port = v
	return nil
}

func Raw(v bool) func(*options) error {
	return func(o *options) error { return o.setRaw(v) }
}
func (o *options) setRaw(v bool) error {
	o.raw = v
	return nil
}
