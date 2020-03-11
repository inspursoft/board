package utils

import (
	"net/http"
)

type FlushResponseWriter struct {
	Writer http.ResponseWriter
}

func (w *FlushResponseWriter) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	if f, ok := w.Writer.(http.Flusher); ok {
		f.Flush()
	}
	return n, err
}
