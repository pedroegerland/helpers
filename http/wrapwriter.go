package http

import "net/http"

type WrapResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
	bodySize    int
}

func (w *WrapResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	// Check after in case there's error handling in the wrapped ResponseWriter.
	if w.wroteHeader {
		return
	}
	w.statusCode = code
	w.wroteHeader = true
}

func (w *WrapResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *WrapResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.bodySize += len(b)
	return w.ResponseWriter.Write(b)
}

func (w *WrapResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
