package middleware

import (
	"net/http"
)

func NOP(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}
