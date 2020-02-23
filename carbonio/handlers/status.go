package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/tabulate/tabulate"
)

const statusTmpl = "html/status.tmpl"

func init() {
	mustTemplate(statusTmpl)
}

func (h *Handlers) StatusCommand(w io.Writer) {
	str, err := status(h.device, h.opts.raw)
	if err != nil {
		helpers.Exit(fmt.Sprintf("error gathering status information; %s", err))
	}

	_, err = io.WriteString(w, str)
	if err != nil {
		helpers.Exit(fmt.Sprintf("error writing status information; %s", err))
	}
}

func (h *Handlers) StatusHandler(w http.ResponseWriter, r *http.Request) {
	stts := http.StatusOK
	buf := &bytes.Buffer{}

	str, err := status(h.device, h.opts.raw)
	if err != nil {
		stts = http.StatusInternalServerError
		w.WriteHeader(stts)
		log.Printf("%s", err)
	}

	if stts == http.StatusOK {
		data := struct {
			Title    string
			Contents string
		}{
			Title:    "Status",
			Contents: str,
		}
		err = tmpls[statusTmpl].Execute(io.Writer(buf), data)
		if err != nil {
			stts = http.StatusInternalServerError
			w.WriteHeader(stts)
			log.Printf("error executing template; %s", err)
		}
	}

	l := 0
	if stts == http.StatusOK {
		l, err = w.Write(buf.Bytes())
		if err != nil {
			stts = http.StatusInternalServerError
			w.WriteHeader(stts)
		}
	}

	helpers.CommonLogFormat(r, stts, l)
}

func status(device devices.Device, asRaw bool) (string, error) {
	lines := []string{"LED STATUS"}

	for _, led := range []*leds.LED{
		device.LEDs().Power(),
		device.LEDs().Status(),
		device.LEDs().Mute(),
	} {
		state, err := led.State()
		if err != nil {
			log.Printf("error reading %s led state; %s", led.Name(), err)
		}
		lines = append(lines, fmt.Sprintf("%s %s", led.Name(), state))
	}

	tbl, err := tabulate.NewTable()
	if err != nil {
		return "", fmt.Errorf("failure creating table; %s", err)
	}
	tbl.Split(lines, ifs, -1)
	rndr := &tabulate.PlainRenderer{}
	rndr.SetOFS(ofs)

	return rndr.Render(tbl), nil
}
