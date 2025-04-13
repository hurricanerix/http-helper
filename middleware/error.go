package middleware

import (
	"fmt"
	"net/http"
	"os"
)

func Error(h http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				os.Stderr.Write([]byte(fmt.Sprintf("ERROR: %v\n", r)))
			}
		}()
		h.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(fn)
}
