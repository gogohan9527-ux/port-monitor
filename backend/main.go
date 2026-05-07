package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	auth, err := NewAuth()
	if err != nil {
		log.Fatalf("init auth: %v", err)
	}

	scanner := NewScanner()
	hub := NewHub()
	loop := NewScannerLoop(scanner, cfg, hub)

	srv := &Server{
		cfg:     cfg,
		auth:    auth,
		scanner: scanner,
		hub:     hub,
		loop:    loop,
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api", func(api chi.Router) {
		api.Post("/login", srv.handleLogin)

		api.Group(func(p chi.Router) {
			p.Use(auth.Middleware)
			p.Get("/ports", srv.handlePorts)
			p.Get("/settings", srv.handleGetSettings)
			p.Put("/settings", srv.handleUpdateSettings)
			p.Get("/ws", srv.handleWS)
		})
	})

	r.Handle("/*", staticHandler())

	go loop.Run()

	snapshot := cfg.Snapshot()
	httpSrv := &http.Server{
		Addr:              snapshot.ListenAddr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		log.Println("shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpSrv.Shutdown(ctx)
		loop.Stop()
		close(idleConnsClosed)
	}()

	log.Printf("port-monitor listening on http://%s (refresh %dms)", snapshot.ListenAddr, snapshot.IntervalMs)
	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v", err)
	}
	<-idleConnsClosed
}
