package python

import (
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := map[string]struct {
		path     string
		wantCode int
		wantBody []byte
	}{
		"existing file":                 {path: "/allowed.txt", wantCode: http.StatusOK, wantBody: []byte("allowed content")},
		"missing file":                  {path: "/missing.txt", wantCode: http.StatusNotFound, wantBody: []byte("")},
		"directory traversal prevented": {path: "/../secrets.txt", wantCode: http.StatusForbidden, wantBody: []byte("")},
	}

	parentDir := t.TempDir()
	sensitiveFile := filepath.Join(parentDir, "secrets.txt")
	if err := os.WriteFile(sensitiveFile, []byte("secret passwords"), 0644); err != nil {
		t.Fatal(err)
	}

	wwwDir := path.Join(parentDir, "www")
	if err := os.Mkdir(wwwDir, 0755); err != nil {
		t.Fatal(err)
	}
	legitimateFile := filepath.Join(wwwDir, "allowed.txt")
	if err := os.WriteFile(legitimateFile, []byte("allowed content"), 0644); err != nil {
		t.Fatal(err)
	}

	handler := Handler{Directory: wwwDir}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			gotCode := rec.Code
			diffCode := cmp.Diff(tc.wantCode, gotCode)
			if diffCode != "" {
				t.Fatal(diffCode)
			}

			gotBody, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatal(err)
			}
			diffBody := cmp.Diff(tc.wantBody, gotBody)
			if diffBody != "" {
				t.Fatal(diffBody)
			}
		})
	}
}
