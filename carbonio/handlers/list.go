package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/signals"
	"github.com/kward/tabulate/tabulate"
)

const listTmpl = "html/list.tmpl"

func init() {
	register(&ListHandler{})
	mustTemplate(listTmpl)
}

// ListHandler implements a handler for status requests.
type ListHandler struct {
	opts   *options
	device devices.Device
}

func NewListHandler(device devices.Device, opts ...func(*options) error) Handler {
	if device == nil {
		helpers.Exit(fmt.Sprintln("device is uninitialized"))
	}

	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			helpers.Exit(fmt.Sprintf("invalid option; %s", err))
		}
	}
	if err := o.validate(); err != nil {
		helpers.Exit(fmt.Sprintf("failed to validate options"))
	}

	return &ListHandler{
		opts:   o,
		device: device,
	}
}

func (h *ListHandler) ServeCommand(w io.Writer) {
	str, err := h.list()
	if err != nil {
		helpers.Exit(fmt.Sprintf("error gathering list information; %s", err))
	}

	_, err = io.WriteString(w, str)
	if err != nil {
		helpers.Exit(fmt.Sprintf("error writing list information; %s", err))
	}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	buf := &bytes.Buffer{}

	str, err := h.list()
	if err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
		log.Printf("%s", err)
	}

	if status == http.StatusOK {
		data := struct {
			Title    string
			Contents string
		}{
			Title:    "List",
			Contents: str,
		}
		err = tmpls[listTmpl].Execute(io.Writer(buf), data)
		if err != nil {
			status = http.StatusInternalServerError
			w.WriteHeader(status)
			log.Printf("error executing template; %s", err)
		}
	}

	l := 0
	if status == http.StatusOK {
		l, err = w.Write(buf.Bytes())
		if err != nil {
			status = http.StatusInternalServerError
			w.WriteHeader(status)
		}
	}

	helpers.CommonLogFormat(r, status, l)
}

func (h *ListHandler) Name() string { return "list" }

func (h *ListHandler) list() (string, error) {
	lines := []string{"SIGNAL GAIN PAD PHANTOM"}

	for i := uint(1); i <= h.device.NumMicInputs(); i++ {
		in, err := h.device.MicInput(i)
		if err != nil {
			log.Printf("error accessing mic input %d; %s\n", i, err)
			continue
		}

		gainStr := gain(h.opts, in)
		padStr := pad(h.opts, in)
		phantomStr := phantom(h.opts, in)
		lines = append(lines, fmt.Sprintf("input/mic/%d %s %s %s", i, gainStr, padStr, phantomStr))
	}

	tbl, err := tabulate.NewTable()
	if err != nil {
		return "", fmt.Errorf("unable to list settings; %s", err)
	}
	tbl.Split(lines, ifs, -1)
	rndr := &tabulate.PlainRenderer{}
	rndr.SetOFS(ofs)

	return rndr.Render(tbl), nil
}

func gain(opts *options, input *signals.Signal) string {
	if !opts.raw {
		v, err := input.Gain()
		if err != nil {
			return "err"
		}
		return fmt.Sprintf("%d", v) // TODO(2020-02-17): include dB unit.
	}
	v, err := input.GainRaw()
	if err != nil {
		return "err"
	}
	return v
}

var boolStr = map[bool]string{
	true:  "On",
	false: "Off",
}

func pad(opts *options, input *signals.Signal) string {
	if !opts.raw {
		v, err := input.Pad()
		if err != nil {
			return "err"
		}
		return boolStr[v]
	}
	v, err := input.PadRaw()
	if err != nil {
		return "err"
	}
	return v
}

func phantom(opts *options, input *signals.Signal) string {
	if !opts.raw {
		v, err := input.Phantom()
		if err != nil {
			return "err"
		}
		return boolStr[v]
	}
	v, err := input.PhantomRaw()
	if err != nil {
		return "err"
	}
	return v
}
