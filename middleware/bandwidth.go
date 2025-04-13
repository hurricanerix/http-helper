package middleware

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/hurricanerix/http-helper/config"
	"github.com/mxk/go-flowrate/flowrate"
)

const defaultBandwidthBps = -1
const defaultBandwidthJitter = 0

func Bandwidth(h http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		brw := &BlockResponse{ResponseWriter: rw, ResponseBuffer: &bytes.Buffer{}}
		h.ServeHTTP(brw, r)
		rw.WriteHeader(brw.StatusCode)

		mbps := config.IntEnv("HH_BANDWIDTH_BPS", defaultBandwidthBps)
		if mbps < 0 {
			_, err := io.Copy(rw, brw.ResponseBuffer)
			if err != nil {
				panic(err)
			}
			return
		}

		ra := rand.New(rand.NewSource(time.Now().Unix()))
		jitter := config.IntEnv("HH_BANDWIDTH_JITTER", defaultBandwidthJitter)
		adjustedJitter := jitter/2 - ra.Intn(jitter)
		limitedReader := flowrate.NewReader(brw.ResponseBuffer, int64(mbps+adjustedJitter))
		_, err := io.Copy(rw, limitedReader)
		if err != nil {
			panic(err)
		}
	}
	return http.HandlerFunc(fn)
}

type BlockResponse struct {
	StatusCode     int
	ResponseBuffer *bytes.Buffer
	ResponseWriter http.ResponseWriter
}

func (b *BlockResponse) Header() http.Header {
	return b.ResponseWriter.Header()
}

func (b *BlockResponse) Write(data []byte) (int, error) {
	return b.ResponseBuffer.Write(data)
}

func (b *BlockResponse) WriteHeader(statusCode int) {
	b.StatusCode = statusCode
}
