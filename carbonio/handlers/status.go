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

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lines := []string{"LED STATUS"}
	var err error

	if h.device == nil {
		fmt.Println("device is uninitialized")
		return
	}

	lines = append(lines, fmt.Sprintf("Power %s", h.device.LEDs().Power()))
	lines = append(lines, fmt.Sprintf("Status %s", h.device.LEDs().Status()))
	lines = append(lines, fmt.Sprintf("Mute %s", h.device.LEDs().Mute()))

	tbl, err := tabulate.NewTable()
	if err != nil {
		fmt.Printf("unable to determine status; %s", err)
		return
	}
	tbl.Split(lines, ifs, -1)
	rndr := &tabulate.PlainRenderer{}
	rndr.SetOFS(ofs)

	data := struct {
		Title  string
		Status string
	}{
		Title:  "Status",
		Status: rndr.Render(tbl),
	}
	buf := &bytes.Buffer{}
	status := http.StatusOK
	err = htmlTmpls["status"].Execute(io.Writer(buf), data)
	if err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
		log.Printf("error executing template; %s", err)
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
