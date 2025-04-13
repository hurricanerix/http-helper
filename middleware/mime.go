package middleware

import (
	"bytes"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
)

func Mime(h http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		mimeType := "application/octet-stream"
		buf := &bytes.Buffer{}
		crw := &CaptureResponse{ResponseWriter: rw, Tee: buf}
		h.ServeHTTP(crw, r)
		if m := mimetype.Detect(buf.Bytes()); m != nil {
			mimeType = m.String()
		}
		rw.Header().Set("Content-Type", mimeType)
	}
	return http.HandlerFunc(fn)
}
