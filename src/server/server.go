package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dragonrider23/go-boiler/src/common"
	"github.com/dragonrider23/verbose"
)

// Server is the main application webserver
type Server struct {
	e         *common.Environment
	routes    http.Handler
	address   string
	HTTPPort  string
	HTTPSPort string
}

// NewServer will initialize a new server with Environment e and routes.
func NewServer(e *common.Environment, routes http.Handler) *Server {
	serv := &Server{
		e:       e,
		routes:  routes,
		address: e.Config.Webserver.Address,
	}

	serv.HTTPPort = strconv.Itoa(e.Config.Webserver.HTTPPort)
	serv.HTTPSPort = strconv.Itoa(e.Config.Webserver.HTTPSPort)
	return serv
}

// Run starts the webserver.
func (s *Server) Run() {
	s.e.Log.Info("Starting web server...")
	if s.e.Config.Webserver.TLSCertFile == "" || s.e.Config.Webserver.TLSKeyFile == "" {
		s.startHTTP()
		return
	}

	if s.e.Config.Webserver.RedirectHTTPToHTTPS {
		go s.startRedirector()
	}
	s.startHTTPS()
}

func (s *Server) startRedirector() {
	s.e.Log.WithFields(verbose.Fields{
		"address": s.address,
		"port":    s.HTTPPort,
	}).Debug()
	s.e.Log.Critical(http.ListenAndServe(
		s.address+":"+s.HTTPPort,
		http.HandlerFunc(s.redirectToHTTPS),
	))
}

func (s *Server) startHTTP() {
	s.e.Log.WithFields(verbose.Fields{
		"address": s.address,
		"port":    s.HTTPPort,
	}).Debug()
	s.e.Log.Fatal(http.ListenAndServe(
		s.address+":"+s.HTTPPort,
		s.routes,
	))
}

func (s *Server) startHTTPS() {
	s.e.Log.WithFields(verbose.Fields{
		"address": s.address,
		"port":    s.HTTPSPort,
	}).Debug()
	s.e.Log.Fatal(http.ListenAndServeTLS(
		s.address+":"+s.HTTPSPort,
		s.e.Config.Webserver.TLSCertFile,
		s.e.Config.Webserver.TLSKeyFile,
		s.routes,
	))
}

func (s *Server) redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	// Lets not do a split if we don't need to
	if s.HTTPPort == "80" && s.HTTPSPort == "443" {
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		return
	}

	host := strings.Split(r.Host, ":")[0]
	http.Redirect(w, r, "https://"+host+":"+s.HTTPSPort+r.RequestURI, http.StatusMovedPermanently)
}
