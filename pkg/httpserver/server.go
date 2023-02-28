package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

// Server - структура http сервера.
type Server struct {
	server          *http.Server
	notify          chan error
	ShutdownTimeout time.Duration
}

// New - Создание нового http сервера.
func New(handler http.Handler, port string) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         port,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		ShutdownTimeout: defaultShutdownTimeout,
	}

	s.start()

	return s
}

// start - запуск сервера.
func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// NotifyError - ожидание ошибки при запуске HTTP-сервера.
func (s *Server) NotifyError() error {
	err := <-s.notify
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("httpserver.Server - NotifyError: %w", err)
	}
	return nil
}

// Shutdown - останавливает сервер.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
