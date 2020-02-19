package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/tabulate/tabulate"
)

const statusTmpl = "html/status.tmpl"

func init() {
	register(&StatusHandler{})
	mustTemplate(statusTmpl)
}

// StatusHandler implements a handler for status requests.
type StatusHandler struct {
	opts   *options
	device devices.Device
}

func NewStatusHandler(device devices.Device, opts ...func(*options) error) Handler {
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

	return &StatusHandler{
		opts:   o,
		device: device,
	}
}

func (h *StatusHandler) ServeCommand(w io.Writer) {
	str, err := h.status()
	if err != nil {
		helpers.Exit(fmt.Sprintf("error gathering status information; %s", err))
	}

	_, err = io.WriteString(w, str)
	if err != nil {
		helpers.Exit(fmt.Sprintf("error writing status information; %s", err))
	}
}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	buf := &bytes.Buffer{}

	str, err := h.status()
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
			Title:    "Status",
			Contents: str,
		}
		err = tmpls[statusTmpl].Execute(io.Writer(buf), data)
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

func (h *StatusHandler) Name() string { return "status" }

func (h *StatusHandler) status() (string, error) {
	lines := []string{"LED STATUS"}
	lines = append(lines, fmt.Sprintf("Power %s", h.device.LEDs().Power()))
	lines = append(lines, fmt.Sprintf("Status %s", h.device.LEDs().Status()))
	lines = append(lines, fmt.Sprintf("Mute %s", h.device.LEDs().Mute()))

	tbl, err := tabulate.NewTable()
	if err != nil {
		return "", fmt.Errorf("failure creating table; %s", err)
	}
	tbl.Split(lines, ifs, -1)
	rndr := &tabulate.PlainRenderer{}
	rndr.SetOFS(ofs)

	return rndr.Render(tbl), nil
}
