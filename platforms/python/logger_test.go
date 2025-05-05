package python

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestLogger(t *testing.T) {
	tests := map[string]struct {
		reqRemoteHost string
		reqTime       time.Time
		reqMethod     string
		reqPath       string
		resStatusCode int
		want          string
	}{
		"http get listing": {
			reqRemoteHost: "127.0.0.1",
			reqTime:       time.Date(2025, 4, 13, 18, 02, 11, 100, time.Local),
			reqMethod:     "GET",
			reqPath:       "/",
			resStatusCode: http.StatusOK,
			want:          "127.0.0.1 - - [13/Apr/2025 18:02:11] \"GET / HTTP/1.1\" 200 -\n"},
		"http head listing": {
			reqRemoteHost: "127.0.0.1",
			reqTime:       time.Date(2025, 4, 13, 18, 02, 11, 100, time.Local),
			reqMethod:     "HEAD",
			reqPath:       "/",
			resStatusCode: http.StatusOK,
			want:          "127.0.0.1 - - [13/Apr/2025 18:02:11] \"HEAD / HTTP/1.1\" 200 -\n"},
		"http head file": {
			reqRemoteHost: "127.0.0.1",
			reqTime:       time.Date(2025, 4, 13, 18, 02, 11, 100, time.Local),
			reqMethod:     "HEAD",
			reqPath:       "/hello.go",
			resStatusCode: http.StatusOK,
			want:          "127.0.0.1 - - [13/Apr/2025 18:02:11] \"HEAD /hello.go HTTP/1.1\" 200 -\n"},
		"http get file": {
			reqRemoteHost: "127.0.0.1",
			reqTime:       time.Date(2025, 4, 13, 18, 02, 11, 100, time.Local),
			reqMethod:     "GET",
			reqPath:       "/hello.go",
			resStatusCode: http.StatusOK,
			want:          "127.0.0.1 - - [13/Apr/2025 18:02:11] \"GET /hello.go HTTP/1.1\" 200 -\n"},
		"http moved permanently file": {
			reqRemoteHost: "127.0.0.1",
			reqTime:       time.Date(2025, 4, 13, 18, 02, 11, 100, time.Local),
			reqMethod:     "GET",
			reqPath:       "/images",
			resStatusCode: http.StatusMovedPermanently,
			want:          "127.0.0.1 - - [13/Apr/2025 18:02:11] \"GET /images HTTP/1.1\" 301 -\n"},
		"http moved permanently": {
			reqRemoteHost: "127.0.0.1",
			reqTime:       time.Date(2025, 4, 13, 18, 02, 11, 100, time.Local),
			reqMethod:     "GET",
			reqPath:       "/images",
			resStatusCode: http.StatusMovedPermanently,
			want:          "127.0.0.1 - - [13/Apr/2025 18:02:11] \"GET /images HTTP/1.1\" 301 -\n"},
		"http delete file": {
			reqRemoteHost: "127.0.0.1",
			reqTime:       time.Date(2025, 4, 13, 18, 02, 11, 100, time.Local),
			reqMethod:     "DELETE",
			reqPath:       "/hello.go",
			resStatusCode: http.StatusNotImplemented,
			want:          "127.0.0.1 - - [13/Apr/2025 18:02:11] \"DELETE /hello.go HTTP/1.1\" 501 -\n"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			preserveTimeNow := timeNow
			preserveStdout := stdout
			defer func() {
				timeNow = preserveTimeNow
				stdout = preserveStdout
			}()
			timeNow = func() time.Time {
				return tc.reqTime
			}
			got := bytes.Buffer{}
			stdout = &got

			svr := httptest.NewServer(Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resStatusCode)
			})))
			defer svr.Close()

			testURL, err := url.JoinPath(svr.URL, tc.reqPath)
			if err != nil {
				t.Errorf("expected err to be nil got %v", err)
			}

			c := http.Client{}
			req, err := http.NewRequest(tc.reqMethod, testURL, nil)
			if err != nil {
				t.Errorf("expected err to be nil got %v", err)
			}

			res, err := c.Do(req)
			if err != nil {
				t.Errorf("expected err to be nil got %v", err)
			}

			if res.StatusCode != tc.resStatusCode {
				t.Errorf("expected statuscode to be %v got %v", tc.resStatusCode, res.StatusCode)
			}

			diff := cmp.Diff(tc.want, got.String())
			if diff != "" {
				t.Fatalf("%s\n", diff)
			}
		})
	}
}
