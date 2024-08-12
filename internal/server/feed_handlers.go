package server

import (
	"Marketplace/internal/entities"
	"Marketplace/internal/lib"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

func (s *Server) GetSorted(w http.ResponseWriter, r *http.Request) {
	log := s.log.With(
		slog.String("method", "GetSorted"),
		slog.String("path", r.URL.Path),
	)

	ReqFeed := &entities.Filter{}

	if err := render.DecodeJSON(r.Body, ReqFeed); err != nil {
		log.Error("Error decoding request body", err)
		render.JSON(w, r, lib.RespError("Error decoding request body"))
		return
	}

	response, err := s.db.GetSorted(r.Context(), ReqFeed)
	if err != nil {
		log.Error("Error getting sorted data", err)
		render.JSON(w, r, lib.RespError("Error getting sorted data"))
		return
	}

	//fmt.Println(response)

	log.Info("Successfully got sorted data")
	render.JSON(w, r, response)
}
