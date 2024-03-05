package log

import "io"

type StreamHandler struct {
	w io.Writer
}

func (that *StreamHandler) Write(p []byte) (n int, err error) {
	return that.w.Write(p)
}

func (that *StreamHandler) Close() error {
	return nil
}

func NewStreamHandler(w io.Writer) *StreamHandler {
	return &StreamHandler{
		w: w,
	}
}
