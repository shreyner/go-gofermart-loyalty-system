package httpserver

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type HTTPServer struct {
	server http.Server
	errors chan error
	log    *zap.Logger
}

func NewHttpServer(log *zap.Logger, handler http.Handler, port int) *HTTPServer {
	return &HTTPServer{
		log: log,
		server: http.Server{
			Addr:    net.JoinHostPort("", strconv.Itoa(port)),
			Handler: handler,
		},
		errors: make(chan error),
	}
}

func (hs *HTTPServer) Start() {
	go func() {
		hs.log.Info("Http Server listening on ", zap.String("addr", hs.server.Addr))
		hs.errors <- hs.server.ListenAndServe()
		close(hs.errors)
	}()
}

func (hs *HTTPServer) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return hs.server.Shutdown(ctx)
}

func (hs *HTTPServer) Notify() <-chan error {
	return hs.errors
}
