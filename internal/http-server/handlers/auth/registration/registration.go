package registration

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
}

type Response struct {
	resp.Response
	UserId int64 `json:"userId"`
}

type UserSaver interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.registration.New"

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

		userId, err := userSaver.RegisterNewUser(r.Context(), req.Email, req.Password)

		if errors.Is(err, auth.ErrUserExists) {
			log.Warn("user already exists", slog.String("email", req.Email))

			render.JSON(w, r, resp.Error("user already exists"))

			return
		}
		if err != nil {
			log.Error("failed to save user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save user"))

			return
		}

		log.Info("user added", slog.Int64("id", userId))

		responseOK(w, r, userId)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, userId int64) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		UserId:   userId,
	})
}
