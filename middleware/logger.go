package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

var logTemplate *template.Template

const loggingTemplateSrc = `{{.RequestTime.Format "2006-01-02T15:04:05-07:00"}} {{.RequestID}} {{.RequestMethod}} {{.RequestPath}} {{tokenWhenEmpty .ResponseStatusCode}} {{tokenWhenEmpty .ResponseContentLength}} {{.ResponseDuration}}`

func init() {
	funcMap := template.FuncMap{
		"tokenWhenEmpty": tokenWhenEmpty,
	}
	logTemplate = template.Must(template.New("logger").Funcs(funcMap).Parse(loggingTemplateSrc))
}

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

		err := logTemplate.Execute(os.Stdout, data)
		if err != nil {
			os.Stderr.Write([]byte(fmt.Sprintf("error writing log: %v\n", err)))
		}
		os.Stdout.Write([]byte{'\n'})
	}
	return http.HandlerFunc(fn)
}

func tokenWhenEmpty(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
