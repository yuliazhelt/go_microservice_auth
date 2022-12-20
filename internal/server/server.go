package server

import (
	"auth/internal/handlers"
	"auth/internal/config"
	"auth/pkg/infra/logger"
	"auth/pkg/infra/tracer"
	"auth/internal/mymiddleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/go-chi/chi/middleware"
	
	"github.com/riandyrn/otelchi"
	
	"auth/internal/adapters/db/pg"
	"go.uber.org/zap"
	"github.com/go-chi/chi/v5"
)

type HttpServer struct {
    server *http.Server
}

func (s *HttpServer) Start() error {
	router := chi.NewRouter()

	db, err := pg.NewPgDb(context.Background(), "postgresql://database:secret@localhost:5432/database")
	if err != nil {
		log.Fatalf("connection to database failed: %#v", err)
	}
	port := config.GetPort()
	salt := os.Getenv("SALT")

	l, err := logger.GetLogger(true, logger.DSN, "myenv")
	if err != nil {
		log.Fatalf("couldn't initialize logger: %#v", err)
	}

	if err := tracer.InitOtel(); err != nil {
		l.Fatal("OTEL init", zap.Error(err))
	}

	// logger middleware

	router.Use(mymiddleware.LoggerMiddleware(l))

	// recovery middleware

	router.Use(middleware.Recoverer) 

	// tracer middleware

	router.Use(otelchi.Middleware("team10_demo_service", otelchi.WithChiRoutes(router))) // tracer middleware копировать библ функцию не надо,можно оставить так а снизу дописать свой middleware
	router.Use(mymiddleware.TracerMiddleware())

	counter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "team10", Name: "counter", Help: "Endpoints request counter",
	})

	//login handler

	router.Post("/auth/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		counter.Inc()
		handlers.Login(w, r, db, salt)
	})

	// verify handler

	router.Post("/auth/api/v1/verify", func(w http.ResponseWriter, r *http.Request) {
		counter.Inc()
		handlers.Verify(w, r)
	})

    s.server = &http.Server{
        Addr:   fmt.Sprintf(":%d", port),
        Handler:    router,
	}

	// metrics to prometheus

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", config.GetMetricsPort()), nil)
    return s.server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
    return s.server.Shutdown(ctx)
}