package python

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hurricanerix/http-helper/middleware"
)

var timeNow = time.Now
var stdout = io.Writer(os.Stdout)
var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(stdout, nil))
}

type logParams struct {
	RemoteHost         string
	RequestTime        time.Time
	RequestMethod      string
	RequestPath        string
	RequestProto       string
	ResponseStatusCode string
}

func newLogParams(crw *middleware.CaptureResponse, r *http.Request) logParams {
	return logParams{
		RemoteHost:         strings.Split(r.Host, ":")[0],
		RequestMethod:      r.Method,
		RequestPath:        r.URL.Path,
		RequestProto:       r.Proto,
		ResponseStatusCode: strconv.FormatInt(int64(crw.StatusCode), 10),
	}
}

func Logger(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		startTime := timeNow()
		crw := &middleware.CaptureResponse{ResponseWriter: rw}
		next.ServeHTTP(crw, r)

		data := newLogParams(crw, r)
		data.RequestMethod = r.Method
		data.RequestTime = startTime

		logger.Info("python_request",
			"remote_host", data.RemoteHost,
			"time", data.RequestTime,
			"method", data.RequestMethod,
			"path", data.RequestPath,
			"proto", data.RequestProto,
			"status", data.ResponseStatusCode,
		)
	}
	return http.HandlerFunc(fn)
}
