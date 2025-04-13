package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/hurricanerix/http-helper/config"
)

const defaultTimeToFirstByte = 400 * time.Millisecond
const defaultTimeToFirstByte95P = 800 * time.Millisecond
const defaultTimeToFirstByteJitter = 200 * time.Millisecond

func TTFB(h http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		ttfb := getTimeToFirstByte()
		time.Sleep(ttfb)
		h.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}

func getTimeToFirstByte() time.Duration {
	ttfb := config.DurationEnv("HH_TIME_TO_FIRST_BYTE", defaultTimeToFirstByte)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	if r.Float64() > 0.95 {
		ttfb = config.DurationEnv("HH_TIME_TO_FIRST_BYTE_95P", defaultTimeToFirstByte95P)
	}
	jitter := config.DurationEnv("HH_TIME_TO_FIRST_BYTE_JITTER", defaultTimeToFirstByteJitter)
	offset, _ := time.ParseDuration(fmt.Sprintf("%fms", 0.5-r.Float64()))
	adjustedJitter := (jitter / time.Millisecond) * offset
	return ttfb + adjustedJitter
}
