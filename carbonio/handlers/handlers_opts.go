package handlers

type options struct {
	// Local flags.
	raw bool
}

func (o *options) validate() error {
	return nil
}

func Raw(v bool) func(*options) error {
	return func(o *options) error { return o.setRaw(v) }
}
func (o *options) setRaw(v bool) error {
	o.raw = v
	return nil
}
