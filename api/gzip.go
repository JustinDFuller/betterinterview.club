package api

import (
	"compress/gzip"
	"io"
	"net/http"
)

type gzipResponseWriter struct {
	w      http.ResponseWriter
	writer io.Writer
}

func (grw gzipResponseWriter) Header() http.Header {
	return grw.w.Header()
}

func (grw gzipResponseWriter) WriteHeader(statusCode int) {
	grw.w.WriteHeader(statusCode)
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.writer.Write(b)
}

func withGzip(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gw := gzip.NewWriter(w)
		defer gw.Close()

		gzr := gzipResponseWriter{w, gw}
		w.Header().Set("Content-Encoding", "gzip")
		handler(gzr, r)
	}
}
