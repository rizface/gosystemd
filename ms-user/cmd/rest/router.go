package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizface/go-ms-systemd/ms-user/database"
)

type router struct {
	r *chi.Mux
}

type Server struct {
	server http.Server
}

func newRouter() *router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return &router{
		r: r,
	}
}

func NewServer(dbPool *pgxpool.Pool) *Server {
	r := newRouter()
	r.registerRoutes(dbPool)

	port := os.Getenv("APP_PORT")
	if port != "" {
		port = ":8000"
	}

	s := http.Server{
		Handler: r.r,
		Addr:    port,
	}

	return &Server{
		server: s,
	}
}

func (r *router) registerRoutes(dbPool *pgxpool.Pool) {
	userRepo := database.NewUser(dbPool)
	userHandler := newUserHandler(userDeps{
		userRepo: userRepo,
	})

	r.r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello")
	})

	r.r.Post("/users", userHandler.CreateUser())
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
