package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muntasirasif/students-api/internal/config"
)

func main() {
	// fmt.Println("Welcome to Students api")

	/// get configuration
	cfg := config.MustLoad()

	///datebase setup

	/// setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student api"))
	})

	/// setup server
	server := http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	slog.Info("starting server", "addr", cfg.Server.Addr)
	// fmt.Println("Server started at", cfg.Server.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		error := server.ListenAndServe()
		if error != nil {
			slog.Error("error starting server", "error", error)
		}
	}()

	<-done
	slog.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("error server stopping", "error", err)
	}

	slog.Info("server stopped successfully")

}
