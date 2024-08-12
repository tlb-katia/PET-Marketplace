package server

import (
	"Marketplace/internal/entities"
	"Marketplace/internal/lib"
	"Marketplace/internal/utils"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	userReq := &entities.User{}
	log := s.log.With(
		slog.String("method", "CreateUser"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	err := render.DecodeJSON(r.Body, &userReq)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", err)
		render.JSON(w, r, lib.RespError("empty request"))
	} else if err != nil {
		log.Error("request body is invalid", err)
		render.JSON(w, r, lib.RespError("invalid request"))
		return
	}

	hashPassword, err := utils.GetHashPassword(userReq.Password)
	if err != nil {
		log.Error("Could not hash password.", err)
		render.JSON(w, r, lib.RespError("Could not hash password."))
		return
	}

	userReq.Password = hashPassword

	response, err := s.db.CreateUser(r.Context(), userReq)
	if err != nil {
		log.Error("Could not create user.", err)
		render.JSON(w, r, lib.RespError("Could not create user."))
		return
	}

	log.Info("User created successfully", slog.Int("user Id", response.Id))

	render.JSON(w, r, response)

}

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	userReq := &entities.LoginReqUser{}

	log := s.log.With(
		slog.String("method", "LoginUser"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	err := render.DecodeJSON(r.Body, &userReq)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", err)
		render.JSON(w, r, lib.RespError("empty request"))
	} else if err != nil {
		log.Error("request body is invalid", err)
		render.JSON(w, r, lib.RespError("invalid request"))
		return
	}

	response, err := s.db.LoginUser(r.Context(), userReq)
	if err != nil {
		log.Error("Could not log in user.", err)
		render.JSON(w, r, lib.RespError("Could not log in user."))
		return
	}

	ok := utils.CheckPassword(response.Password, userReq.Password)
	if !ok {
		log.Error("The password is incorrect", err)
		render.JSON(w, r, lib.RespError("The password is incorrect"))
		return
	}

	token, err := utils.GenerateToken(response.Id)
	if err != nil {
		log.Error("Token wasn't created", err)
		render.JSON(w, r, lib.RespError("Token wasn't created"))
		return
	}

	log.Info("User logged in successfully", slog.Int("user Id", response.Id))
	render.JSON(w, r, token)
}
