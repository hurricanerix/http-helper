package python

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
)

const directoryListingTemplateSrc = `<!DOCTYPE HTML>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Directory listing for {{.Path}}</title>
</head>
<body>
<h1>Directory listing for {{.Path}}</h1>
<hr>
<ul>{{range $name, $path := .Files}}
	<li><a href="{{$path}}">{{$name}}</a></li>
{{end}}</ul>
<hr>
</body>
</html>`

var listingTemplate *template.Template

func init() {
	var err error
	listingTemplate, err = template.New("directoryListing").Parse(directoryListingTemplateSrc)
	if err != nil {
		panic(err)
	}
}

type Handler struct {
	Directory string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.ReadAll(r.Body)

	p := strings.TrimRight(r.URL.Path, "/")
	target := path.Join(h.Directory, p) // TODO: ABS?

	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		h.serveListing(target, w, r)
		return
	}

	w.Header().Add("Server", "SimpleHTTP/0.6 Python/3.12.3")
	h.serveFile(target, w, r)
}

func (h Handler) serveFile(target string, w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(target)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}

type ListingTemplateData struct {
	Path  string
	Files map[string]string
}

func (h Handler) serveListing(target string, w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(target)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files := make(map[string]string, len(entries))
	for _, e := range entries {
		files[e.Name()] = path.Join(r.URL.Path, e.Name())
	}

	data := ListingTemplateData{
		Path:  r.URL.Path,
		Files: files,
	}

	payload := &bytes.Buffer{}

	err = listingTemplate.Execute(payload, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to execute template: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload.Bytes())
}

/*
< HTTP/1.0 200 OK
< Server: SimpleHTTP/0.6 Python/3.12.3
< Date: Fri, 02 May 2025 17:46:37 GMT
< Content-type: text/html; charset=utf-8
< Content-Length: 322
<
<!DOCTYPE HTML>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Directory listing for /</title>
</head>
<body>
<h1>Directory listing for /</h1>
<hr>
<ul>
<li><a href="hello.go">hello.go</a></li>
<li><a href="images/">images/</a></li>
<li><a href="WaterBottle.glb">WaterBottle.glb</a></li>
</ul>
<hr>
</body>
</html>
* Closing connection
*/

/*
< HTTP/1.0 200 OK
< Server: SimpleHTTP/0.6 Python/3.12.3
< Date: Fri, 02 May 2025 17:47:48 GMT
< Content-type: application/octet-stream
< Content-Length: 74
< Last-Modified: Sun, 13 Apr 2025 02:05:23 GMT
<
package main

import "fmt"

func main() {
	fmt.Println("Hello, 世界")
}
* Closing connection
*/

/*
$ curl -d -XDELETE localhost:8000/hello.go
<!DOCTYPE HTML>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>Error response</title>
    </head>
    <body>
        <h1>Error response</h1>
        <p>Error code: 501</p>
        <p>Message: Unsupported method ('POST').</p>
        <p>Error code explanation: 501 - Server does not support this operation.</p>
    </body>
</html>

*/
