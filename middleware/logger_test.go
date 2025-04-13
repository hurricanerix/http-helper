package middleware

import (
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	//tests := map[string]struct {
	//	inputTemplate string
	//	inputRequest  http.Request
	//	want          string
	//}{
	//	"simple": {inputTemplate: "", want: ""},
	//}
	//
	//logger := NewLogger("")
	//
	//svr := httptest.NewServer(logger.Wrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(http.StatusOK)
	//})))
	//defer svr.Close()
	//
	//for name, tc := range tests {
	//	t.Run(name, func(t *testing.T) {
	//		c := http.Client{}
	//		res, err := c.Get(svr.URL)
	//		if err != nil {
	//			t.Errorf("expected err to be nil got %v", err)
	//		}
	//		diff := cmp.Diff(tc.want, got)
	//		if diff != "" {
	//			t.Fatalf(diff)
	//		}
	//	})
	//}
}
