package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/tabulate/tabulate"
)

// Status implements a handler for status requests.
type StatusHandler struct {
	device devices.Device
}

// Ensure the Handler interface is implemented.
var _ Handler = new(StatusHandler)

func NewStatusHandler(device devices.Device) Handler {
	if device == nil {
		helpers.Exit(fmt.Sprintln("device is uninitialized"))
	}

	return &StatusHandler{device: device}
}

func init() {
	tmpl, err := template.New("status").Parse(statusTmpl)
	if err != nil {
		helpers.Exit(fmt.Sprintf("error parsing status template; %s", err))
	}
	htmlTmpls["status"] = tmpl
}

var statusTmpl = `
<html>
<head>
<title>{{.Title}}</title>
</head>

<body>
<pre>
{{.Status}}
</pre>
</body>
</html>
`

func (h *StatusHandler) ServeCommand(w io.Writer) {
	str, err := h.status()
	if err != nil {
		helpers.Exit(fmt.Sprintf("error gathering status; %s", err))
	}

	_, err = io.WriteString(w, str)
	if err != nil {
		helpers.Exit(fmt.Sprintf("error writing status; %s", err))
	}
}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	buf := &bytes.Buffer{}

	str, err := h.status()
	if err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
		log.Printf("error executing template; %s", err)
	}

	if status == http.StatusOK {
		data := struct {
			Title  string
			Status string
		}{
			Title:  "Status",
			Status: str,
		}
		err = htmlTmpls["status"].Execute(io.Writer(buf), data)
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
	var err error

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
