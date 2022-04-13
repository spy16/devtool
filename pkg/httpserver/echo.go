package httpserver

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// Echo returns a request echo handler.
func Echo() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			_, _ = wr.Write([]byte(fmt.Sprintf("failed to dump: ")))
		} else {
			wr.WriteHeader(http.StatusOK)
			_, _ = wr.Write(dump)
		}
	})
}

// TeaPot server that does nothing.
func TeaPot() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		wr.WriteHeader(http.StatusTeapot)
	})
}
