package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

type LogParams struct {
	// Common Log Format
	RemoteHost            string
	ClientIdentity        string
	UserID                string
	RequestTime           time.Time
	RequestMethod         string
	RequestPath           string
	RequestProto          string
	ResponseStatusCode    string
	ResponseContentLength string

	// Extra
	RequestID           string
	ResponseDuration    time.Duration
	ResponseContentType string
	ResponseETag        string
}

func NewLogParams(crw *CaptureResponse, r *http.Request) LogParams {
	return LogParams{

		ClientIdentity:        "",
		UserID:                "",
		RequestMethod:         r.Method,
		RequestPath:           r.URL.Path,
		RequestProto:          r.Proto,
		ResponseStatusCode:    strconv.FormatInt(int64(crw.StatusCode), 10),
		ResponseContentLength: strconv.FormatInt(crw.BytesWritten, 10),
		RequestID:             crw.Header().Get("X-Request-Id"),
		ResponseContentType:   crw.Header().Get("Content-Type"),
		ResponseETag:          crw.Header().Get("ETag"),
	}
}

func Logger(h http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		crw := &CaptureResponse{ResponseWriter: rw}
		h.ServeHTTP(crw, r)
		elapsedTime := time.Since(startTime)

		data := NewLogParams(crw, r)
		data.RequestTime = startTime
		data.ResponseDuration = elapsedTime

		logger.Info("request",
			"time", data.RequestTime,
			"request_id", data.RequestID,
			"method", data.RequestMethod,
			"path", data.RequestPath,
			"status", data.ResponseStatusCode,
			"bytes", data.ResponseContentLength,
			"duration", data.ResponseDuration,
			"content_type", data.ResponseContentType,
			"etag", data.ResponseETag,
		)
	}
	return http.HandlerFunc(fn)
}
