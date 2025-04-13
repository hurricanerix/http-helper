package middleware

import (
	"io"
	"net/http"
)

type CaptureResponse struct {
	StatusCode     int
	statusWritten  bool
	BytesWritten   int64
	ResponseWriter http.ResponseWriter
	Tee            io.Writer
}

func (c *CaptureResponse) GetHeader(value string) string {
	return c.ResponseWriter.Header().Get(value)
}

func (c *CaptureResponse) Header() http.Header {
	return c.ResponseWriter.Header()
}

func (c *CaptureResponse) Write(data []byte) (int, error) {
	n, err := c.ResponseWriter.Write(data)
	c.BytesWritten += int64(n)
	if c.Tee != nil {
		c.Tee.Write(data)
	}
	return n, err
}

func (c *CaptureResponse) WriteHeader(statusCode int) {
	c.ResponseWriter.WriteHeader(statusCode)
	if c.statusWritten {
		return
	}
	c.StatusCode = statusCode
}
