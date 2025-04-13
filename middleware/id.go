package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func RequestID(h http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		id := strings.ReplaceAll(uuid.New().String(), "-", "")
		rw.Header().Add("X-Request-ID", id)
		h.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}
