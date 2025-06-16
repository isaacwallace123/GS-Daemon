package routes

import (
	"net/http"
	"strings"

	"Daemon/internal/api/handlers"
)

func RegisterRoutes(mux *http.ServeMux, h *handlers.ContainerHandler) {
	mux.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Create(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/containers/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/containers/")
		segments := strings.Split(path, "/")
		name := segments[0]

		switch {
		case r.Method == http.MethodGet && len(segments) == 1:
			h.Get(w, r, name)
		case r.Method == http.MethodDelete && len(segments) == 1:
			h.Delete(w, r, name)
		case r.Method == http.MethodPost && len(segments) == 2 && segments[1] == "start":
			h.Start(w, r, name)
		case r.Method == http.MethodPost && len(segments) == 2 && segments[1] == "stop":
			h.Stop(w, r, name)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

}
