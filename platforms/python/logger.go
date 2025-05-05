package python

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/hurricanerix/http-helper/middleware"
)

var timeNow = time.Now
var stdout = io.Writer(os.Stdout)

var logTemplate *template.Template

const loggingTemplateSrc = `{{.RemoteHost}} - - [{{.RequestTime.Format "02/Jan/2006 15:04:05"}}] "{{.RequestMethod}} {{.RequestPath}} {{.RequestProto}}" {{tokenWhenEmpty .ResponseStatusCode}} -`

func init() {
	funcMap := template.FuncMap{
		"tokenWhenEmpty": tokenWhenEmpty,
	}
	logTemplate = template.Must(template.New("logger").Funcs(funcMap).Parse(loggingTemplateSrc))
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

		f := bufio.NewWriter(stdout)
		err := logTemplate.Execute(f, data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing log: %v\n", err)
		}
		f.Write([]byte{'\n'})
		f.Flush()
	}
	return http.HandlerFunc(fn)
}

func tokenWhenEmpty(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
