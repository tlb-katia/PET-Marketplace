package server

import (
	"Marketplace/internal/entities"
	"Marketplace/internal/lib"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (ad *Server) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	advertReq := &entities.Advert{}
	log := ad.log.With(
		slog.String("method", "CreateAdvert"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	err := render.DecodeJSON(r.Body, &advertReq)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", err)
		render.JSON(w, r, lib.RespError("empty request"))
	} else if err != nil {
		log.Error("request body is invalid", err)
		render.JSON(w, r, lib.RespError("invalid request"))
		return
	}

	advertReq.UserId = r.Context().Value("userID").(int)
	advertReq.Datetime = time.Now()
	advertReq.ByThisUser = true

	response, err := ad.db.CreateAdvert(r.Context(), advertReq)
	if err != nil {
		log.Error("Could not create advert.", err)
		render.JSON(w, r, lib.RespError("Could not create advert."))
		return
	}

	log.Info("Advert created successfully", slog.Int("advert Id", advertReq.Id))
	render.JSON(w, r, response)
}

func (ad *Server) GetAdvert(w http.ResponseWriter, r *http.Request) {
	log := ad.log.With(
		slog.String("method", "GetAdvert"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	advertIDStr := chi.URLParam(r, "id")
	advertID, err := strconv.Atoi(advertIDStr)
	if err != nil {
		log.Error("no id in the url", err)
		render.JSON(w, r, lib.RespError("no id in the url"))
	}

	response, err := ad.db.GetAdvert(r.Context(), advertID)
	if err != nil {
		log.Error("Could not get advert", err)
		render.JSON(w, r, lib.RespError("Could not get advert"))
		return
	}
	log.Info("Advert found successfully", slog.Int("advert Id", response.Id))
	render.JSON(w, r, response)
}

func (ad *Server) UpdateAdvert(w http.ResponseWriter, r *http.Request) {
	log := ad.log.With(
		slog.String("method", "UpdateAdvert"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	advertIDStr := chi.URLParam(r, "id")
	advertID, err := strconv.Atoi(advertIDStr)
	if err != nil {
		log.Error("no id in the url", err)
		render.JSON(w, r, lib.RespError("no id in the url"))
		return
	}

	advertReq := &entities.Advert{}
	if err := render.DecodeJSON(r.Body, advertReq); err != nil {
		log.Error("request body is invalid", err)
		render.JSON(w, r, lib.RespError("request body is invalid"))
	}

	if err, _ = ad.checkAdOwner(r, advertID); err != nil {
		log.Error("problems with checking the owner", err)
		render.JSON(w, r, lib.RespError("problems with checking the owner"))
		return
	}
	advertReq.Id = advertID
	advertReq.ByThisUser = true

	response, err := ad.db.UpdateAdvert(r.Context(), advertReq)
	if err != nil {
		log.Error("couldn't update the ad")
		render.JSON(w, r, lib.RespError("couldn't update the ad"))
		return
	}

	log.Info("The ad is updated successfully", slog.Int("advert Id", response.Id))
	render.JSON(w, r, response)
}

func (ad *Server) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	log := ad.log.With(
		slog.String("method", "UpdateAdvert"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	adIdStr := chi.URLParam(r, "id")
	adId, err := strconv.Atoi(adIdStr)
	if err != nil {
		log.Error("no id in the url", err)
		render.JSON(w, r, lib.RespError("no id in the url"))
		return
	}

	err, adToDelete := ad.checkAdOwner(r, adId)
	if err != nil {
		log.Error("problems with checking the owner", err)
		render.JSON(w, r, lib.RespError("problems with checking the owner"))
		return
	}

	err = ad.db.DeleteAdvert(r.Context(), adId)
	if err != nil {
		log.Error("couldn't delete the ad", err)
		render.JSON(w, r, lib.RespError("couldn't delete the ad"))
		return
	}

	log.Info("The ad is deleted successfully", slog.Int("ad Id", adId))
	render.JSON(w, r, adToDelete)
}
