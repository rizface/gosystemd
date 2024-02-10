package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rizface/go-ms-systemd/ms-user/user"
)

type userDeps struct {
	userRepo user.UserRepo
}

type userHandler struct {
	userDeps userDeps
}

func newUserHandler(userDeps userDeps) *userHandler {
	return &userHandler{
		userDeps: userDeps,
	}
}

func (u *userHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var usr user.User

		if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
			slog.ErrorContext(r.Context(), "failed parse payload: %w", err)
			returnResponse(w, Response{
				Code: http.StatusBadRequest,
				Info: "parse error",
			})
			return
		}

		usr, err := user.CreateUser(r.Context(), u.userDeps.userRepo, usr)
		if errors.Is(err, user.ErrUsernameAlreadyTaken) {
			returnResponse(w, Response{
				Code: http.StatusConflict,
				Info: err.Error(),
			})
			return
		}

		if err != nil {
			slog.ErrorContext(r.Context(), "failed create user: %w", err)
			returnResponse(w, Response{
				Code: http.StatusInternalServerError,
				Info: "internal server error",
			})
			return
		}

		returnResponse(w, Response{
			Code: http.StatusOK,
			Info: "success",
			Data: usr,
		})
	}
}
