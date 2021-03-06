package httpserver

import (
	"context"
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"go.uber.org/zap"
)

// HTTPServer is http message delivery
type HTTPServer struct {
	logger *zap.Logger
	addr   string
}

// New created new HTTPServer
func New(logger *zap.Logger, ip string, port int) *HTTPServer {
	return &HTTPServer{
		logger: logger,
		addr:   fmt.Sprintf("%s:%d", ip, port),
	}
}

// AddHandler added handler to server
func (s *HTTPServer) AddHandler(pattern string, handler domain.Handler) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			s.logger.Error("error read data from html request", zap.Error(err))
			writer.WriteHeader(400)
			writer.Write([]byte(err.Error()))
			return
		}

		response, err := handler(&domain.Event{
			Heading: string(data),
		})
		if err != nil {
			s.logger.Error("error handle request", zap.Error(err))
			writer.WriteHeader(400)
			writer.Write([]byte(err.Error()))
			return
		}

		s.logger.Debug("handle ran correct", zap.Error(err))
		writer.WriteHeader(200)
		writer.Write([]byte(response))
		return
	})
}

// Run server (blocked)
func (s *HTTPServer) Run(ctx context.Context) error {
	s.logger.Debug("http server starting")

	srv := &http.Server{Addr: s.addr}

	go func(ctx context.Context, server *http.Server) {
		<-ctx.Done()
		err := srv.Shutdown(ctx)
		if err != nil {
			s.logger.Error("http cannot stopping", zap.Error(err))
		}
	}(ctx, srv)

	return srv.ListenAndServe()
}
