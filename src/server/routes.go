package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/dragonrider23/go-boiler/src/common"
	mid "github.com/dragonrider23/go-boiler/src/server/middleware"

	"github.com/gorilla/mux"
)

// LoadRoutes collects all application routes and registers them with a router.
func LoadRoutes(e *common.Environment) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Page routes
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.HandleFunc("/", rootHandler)
	r.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	// Development Routes
	if e.Dev {
		// Add routes for profiler
		s := r.PathPrefix("/debug").Subrouter()
		s.HandleFunc("/pprof/", pprof.Index)
		s.HandleFunc("/pprof/cmdline", pprof.Cmdline)
		s.HandleFunc("/pprof/profile", pprof.Profile)
		s.HandleFunc("/pprof/symbol", pprof.Symbol)
		s.HandleFunc("/pprof/trace", pprof.Trace)
		// Manually add support for paths linked to by index page at /debug/pprof/
		s.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
		s.Handle("/pprof/heap", pprof.Handler("heap"))
		s.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
		s.Handle("/pprof/block", pprof.Handler("block"))
		e.Log.Debug("Profiling enabled")
	}

	h := mid.Logging(e, r) // Logging

	return h
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
