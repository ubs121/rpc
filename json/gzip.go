package json

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipResponseWriter is used for gzipped response.
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// makes gzip handler
func MakeGzipHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// gzip
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()

			handler.ServeHTTP(GzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
		} else {
			handler.ServeHTTP(w, r)
		}

	})
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
