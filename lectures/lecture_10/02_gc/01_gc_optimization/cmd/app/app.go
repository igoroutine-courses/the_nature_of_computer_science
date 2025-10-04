package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

const addr = "0.0.0.0:8080"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	s := &server{
		logger: logger,
		pool: &sync.Pool{
			New: func() any {
				return make([]byte, 10<<20)
			},
		},
	}

	debug.SetGCPercent(100)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/action", s.endPoint)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		defer cancel()
		if err := server.ListenAndServe(); err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"server serve error",
				slog.Any("error", err),
			)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := server.Shutdown(shutdownCtx)

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"server shutdown error",
			slog.Any("error", err),
		)
		return
	}
}

type server struct {
	logger *slog.Logger
	pool   *sync.Pool
}

func (s *server) endPoint(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("request")

	// We can use sync.Pool here
	//buffer := s.pool.Get().([]byte)
	buffer := make([]byte, 1<<20)
	//defer s.pool.Put(buffer)
	w.Write(buffer[:10])
	w.WriteHeader(http.StatusOK)
}
