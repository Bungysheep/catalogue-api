package rest

import (
	"net/http"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/configs"
	"github.com/bungysheep/catalogue-api/pkg/protocols/rest/routes"
)

// Server type
type Server struct {
	*http.Server
}

// NewRestServer creates new rest server
func NewRestServer() *Server {
	return &Server{}
}

// RunServer runs rest server
func (s *Server) RunServer() error {

	s.Server = &http.Server{
		Addr:         ":" + configs.PORT,
		Handler:      routes.APIV1RouteHandler(),
		ReadTimeout:  configs.READTIMEOUT * time.Second,
		WriteTimeout: configs.WRITETIMEOUT * time.Second,
	}

	return s.Server.ListenAndServe()
}
