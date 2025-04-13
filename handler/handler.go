package handler

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

type File struct {
	Directory string
}

func (h File) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	h.serveFile(target, w, r)
}

func (h File) serveFile(target string, w http.ResponseWriter, r *http.Request) {
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

func (h File) serveListing(target string, w http.ResponseWriter, r *http.Request) {
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
