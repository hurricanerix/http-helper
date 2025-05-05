package s3

import (
	"net/http"
)

type Handler struct {
	Directory string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("TODO: implement a mock S3 API which read and writes to the filesystem."))
}
