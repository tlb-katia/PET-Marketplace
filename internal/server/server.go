package server

import (
	"Marketplace/config"
	"Marketplace/internal/entities"
	"Marketplace/internal/server/mw"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type AdvertProvider interface {
	CreateAdvert(ctx context.Context, ad *entities.Advert) (*entities.Advert, error)
	GetAdvert(ctx context.Context, id int) (*entities.Advert, error)
	UpdateAdvert(ctx context.Context, ad *entities.Advert) (*entities.Advert, error)
	DeleteAdvert(ctx context.Context, id int) error
	GetSorted(ctx context.Context, filter *entities.Filter) (*entities.AdvList, error)
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	LoginUser(ctx context.Context, user *entities.LoginReqUser) (*entities.User, error)
	GetUserByLogin(login string) (*entities.User, error)
	CheckUserExists(login string) bool
}

type Server struct {
	db AdvertProvider
	//redisDB
	router *chi.Mux
	log    *slog.Logger
}

func NewServer(db AdvertProvider, router *chi.Mux, log *slog.Logger) *Server {
	return &Server{
		db:     db,
		router: router,
		log:    log,
	}
}

func (s *Server) Run(config *config.Config) {

	s.router.Use(middleware.RequestID)
	//s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.URLFormat)

	s.router.Route("/auth", func(r chi.Router) {
		r.Post("/", s.CreateUser)
		r.Get("/{id}", s.LoginUser)
	})

	s.router.Route("/api", func(r chi.Router) {
		r.With(mw.Auth).Post("/", s.CreateAdvert)
		r.Get("/{id}", s.GetAdvert)
		r.With(mw.Auth).Put("/{id}", s.UpdateAdvert)
		r.With(mw.Auth).Delete("/{id}", s.DeleteAdvert)
	})
	s.router.Get("/feed", s.GetSorted)

	srv := http.Server{
		Addr:    config.HTTPServerPort,
		Handler: s.router,
	}

	//go func() {
	if err := srv.ListenAndServe(); err != nil {
		s.log.Error("failed to start server")
	}
	//}()
}
