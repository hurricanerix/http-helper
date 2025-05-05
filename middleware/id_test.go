package middleware

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestRequestIDMiddleware(t *testing.T) {
	uuid.SetRand(rand.New(rand.NewSource(1)))

	svr := httptest.NewServer(RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	defer svr.Close()

	c := http.Client{}
	res, err := c.Get(svr.URL)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	want := "52fdfc072182454f963f5f0f9a621d72"
	id := res.Header.Get("X-Request-Id")
	if id != want {
		t.Errorf("want %v, got %v", want, id)
	}
}
