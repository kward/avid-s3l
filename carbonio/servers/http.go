package servers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/handlers"
	"github.com/kward/avid-s3l/carbonio/helpers"
)

var (
	htmlTmpls map[string]*template.Template = map[string]*template.Template{}
	device    devices.Device
)

func HttpServer(port int, dev devices.Device) {
	device = dev
	log.SetFlags(0)

	http.HandleFunc("/", rootHandler)

	statusHandler := handlers.NewStatusHandler(device)
	http.Handle("/status", statusHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("Hello, world!\n"))
	buf.WriteString(fmt.Sprintf("r.URL.Path: %s\n", r.URL.Path))

	status := http.StatusOK
	l, err := w.Write(buf.Bytes())
	if err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
	}
	helpers.CommonLogFormat(r, status, l)
}
