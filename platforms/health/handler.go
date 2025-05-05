package health

import (
	"net/http"
)

// Handler responsible for handeling readiness probe requests indicating wheather the service
// is ready to receive requests.
type Handler struct {
}

// ServeHTTP returns a 200 OK on a GET request and 405 Method Not Supported on all other
// requests.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Content-Length", "0")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//resp := Response{
	//	BuildInfo: BuildInfo{
	//		VCSRevision: buildinfo.VCSRevision,
	//		VCSModified: buildinfo.VCSModified,
	//	},
	//}

	//b, _ := json.Marshal(&resp)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//n, err := w.Write(b)
	//if err != nil {
	//	panic(fmt.Sprintf("could not write response body: wrote %d bytes of %d: %s", n, len(b), err))
	//}
}
