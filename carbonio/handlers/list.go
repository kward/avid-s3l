package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/tabulate/tabulate"
)

const listTmpl = "html/list.tmpl"

func init() {
	mustTemplate(listTmpl)
}

func (h *Handlers) ListCommand(w io.Writer) {
	str, err := list(h.device, h.opts.raw)
	if err != nil {
		// TODO(2020-02-23) Add an error to `h` instead of exiting here.
		helpers.Exit(fmt.Sprintf("error gathering list information; %s", err))
	}

	_, err = io.WriteString(w, str)
	if err != nil {
		helpers.Exit(fmt.Sprintf("error writing list information; %s", err))
	}
}

func (h *Handlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	stts := http.StatusOK

	data := struct {
		Title string
		Host  net.IP
		Port  int
	}{
		Title: "List",
		Host:  h.device.IP(),
		Port:  h.opts.port,
	}
	if err := tmpls[listTmpl].Execute(io.Writer(buf), data); err != nil {
		stts = http.StatusInternalServerError
		w.WriteHeader(stts)
		log.Printf("error executing template; %s", err)
	}

	l := 0
	if stts == http.StatusOK {
		var err error
		l, err = w.Write(buf.Bytes())
		if err != nil {
			stts = http.StatusInternalServerError
			w.WriteHeader(stts)
		}
	}

	helpers.CommonLogFormat(r, stts, l)
}

func (h *Handlers) ListQueryHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	stts := http.StatusOK

	page, err := list(h.device, h.opts.raw)
	if err != nil {
		stts = http.StatusInternalServerError
		w.WriteHeader(stts)
		log.Printf("%s", err)
	}

	if stts == http.StatusOK {
		buf.WriteString(page)
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

var boolToStr = map[bool]string{
	true:  "On",
	false: "Off",
}

func list(device devices.Device, asRaw bool) (string, error) {
	lines := []string{"SIGNAL GAIN PAD PHANTOM"}

	for i := 1; i <= device.NumMicInputs(); i++ {
		s, err := device.MicInput(i)
		if err != nil {
			log.Printf("error accessing mic input %d; %s", i, err)
			continue
		}

		if asRaw {
			lines = append(lines, fmt.Sprintf("input/mic/%d %q %q %q", i,
				s.Gain().Raw(), s.Pad().Raw(), s.Phantom().Raw()))
			continue
		}
		gain, err := s.Gain().Value()
		if err != nil {
			log.Printf("error reading mic input %d gain; %s", i, err)
		}
		pad, err := s.Pad().IsEnabled()
		if err != nil {
			log.Printf("error read mic input %d pad; %s", i, err)
		}
		phantom, err := s.Phantom().IsEnabled()
		if err != nil {
			log.Printf("error read mic input %d phantom; %s", i, err)
		}
		lines = append(lines, fmt.Sprintf("input/mic/%d %d %s %s", i,
			gain, boolToStr[pad], boolToStr[phantom]))
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
