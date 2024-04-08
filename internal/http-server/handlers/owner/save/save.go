package save

import (
	"effective_mobile_test/internal/lib/api/response"
	"effective_mobile_test/internal/lib/logger/sl"
	"effective_mobile_test/internal/storage/postgres"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	postgres.Owner
}

type Response struct {
	response.Response
	OwnerId int `json:"owner_id"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=OwnerSaver
type OwnerSaver interface {
	SaveOwner(owner postgres.Owner) (int, error)
}

// @Summary		Save a new owner
// @Description	Save a new owner by name, surname, patronymic
// @Tags			Owner
// @Accept			json
// @Produce		json
// @Param			name		body		string	true	"Name"
// @Param			surname		body		string	true	"Surname"
// @Param			patronymic	body		string	true	"Patronymic"
// @Success		200			{object}	Response
// @Failure		400			{object}	response.Response
// @Router			/owner/save [post]
func New(log *slog.Logger, ownerSaver OwnerSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.owner.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if ok, field, msg := validateRequest(req); !ok {
			log.Error("invalid request", field)

			render.JSON(w, r, response.Error(msg))

			return
		}

		ownerId, err := ownerSaver.SaveOwner(req.Owner)
		if err != nil {
			log.Error("failed to save owner", sl.Err(err))

			render.JSON(w, r, response.Error("failed to save owner"))

			return
		}

		render.JSON(w, r, Response{
			response.OK(),
			ownerId,
		})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if req.Name == "" {
		return false, slog.String("field", "name"), "field name is not valid"
	}
	if req.Surname == "" {
		return false, slog.String("field", "surname"), "field surname is not valid"
	}
	if req.Patronymic == "" {
		return false, slog.String("field", "patronymic"), "field patronymic is not valid"
	}
	return true, slog.Attr{}, ""
}
