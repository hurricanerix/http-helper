package middleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
)

func ETag(h http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}
		crw := &CaptureResponse{ResponseWriter: rw, Tee: buf}
		h.ServeHTTP(crw, r)
		h := md5.New()
		h.Write(buf.Bytes())
		hash := fmt.Sprintf("%x", h.Sum(nil))
		rw.Header().Set("ETag", hash)
	}
	return http.HandlerFunc(fn)
}
