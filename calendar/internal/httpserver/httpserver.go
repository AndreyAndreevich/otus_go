package httpserver

import (
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"go.uber.org/zap"
)

// HttpServer is http message delivery
type HttpServer struct {
	logger *zap.Logger
	addr   string
}

// New created new HttpServer
func New(logger *zap.Logger, ip string, port int) *HttpServer {
	return &HttpServer{
		logger: logger,
		addr:   fmt.Sprintf("%s:%d", ip, port),
	}
}

// AddHandler added handler to server
func (s *HttpServer) AddHandler(pattern string, handler domain.Handler) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			s.logger.Error("error read data from html request", zap.Error(err))
			writer.WriteHeader(400)
			writer.Write([]byte(err.Error()))
			return
		}

		response, err := handler(domain.EventData(data))
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
func (s *HttpServer) Run() error {
	return http.ListenAndServe(s.addr, nil)
}
