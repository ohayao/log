package log

import (
	"os"
)

type terminalHandler struct {
	w *os.File
}

func (t *terminalHandler) Write(b []byte) (n int, err error) {
	return t.w.Write(b)
}

func (t *terminalHandler) Close() (err error) {
	return t.w.Close()
}

func newTerminalHandler(file *os.File) *terminalHandler {
	if file == nil {
		file = os.Stdout
	}
	return &terminalHandler{w: file}
}
