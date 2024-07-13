package login

import (
	"context"
	"errors"
	"io"
	"net/http"

	"log/slog"

	resp "driver_backend/internal/lib/api/response"
	"driver_backend/internal/lib/logger/sl"
	"driver_backend/internal/services/auth"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	AppId    int    `json:"appId"`
}

type Response struct {
	resp.Response
	Token string `json:"token"`
}

type UserProvider interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
}

func New(log *slog.Logger, userProvider UserProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.login.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		token, err := userProvider.Login(r.Context(), req.Email, req.Password, int(req.AppId))

		if errors.Is(err, auth.ErrInvalidCredentials) {
			log.Warn("invalid email or password", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid email or password"))

			return
		}
		if err != nil {
			log.Warn("failed to login", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to login"))

			return
		}

		responseOK(w, r, token)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, token string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Token:    token,
	})
}
