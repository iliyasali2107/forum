package server

import (
	"context"
	"net/http"
)

type Server struct {
	Srv *http.Server
}

// func NewServer(cfg *config.Config, handler *delivery.Controller) *Server {
// 	mux := http.NewServeMux()
// 	handler.InitRoutes(mux)
// 	return &Server{
// 		Srv: &http.Server{
// 			Addr:           cfg.Http.Addr,
// 			Controller:        mux,
// 			ReadTimeout:    1,
// 			WriteTimeout:   1,
// 			MaxHeaderBytes: cfg.Http.MaxHeaderBytes << 20,
// 		},
// 	}
// }

func (srv *Server) Run() error {
	return srv.Srv.ListenAndServe()
}

func (srv *Server) ShutDown(ctx context.Context) error {
	return srv.Srv.Shutdown(ctx)
}
