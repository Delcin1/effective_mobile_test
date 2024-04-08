package update

import (
	"effective_mobile_test/internal/lib/api/response"
	"effective_mobile_test/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	OwnerId    int     `json:"ownerId"`
	Name       *string `json:"name"`
	Surname    *string `json:"surname"`
	Patronymic *string `json:"patronymic"`
}

type Response struct {
	response.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=OwnerUpdater
type OwnerUpdater interface {
	UpdateOwnerName(ownerID int, newName string) error
	UpdateOwnerSurname(ownerID int, newSurname string) error
	UpdateOwnerPatronymic(ownerID int, newPatronymic string) error
}

// @Summary		Update owner
// @Description	Update owner by ownerId and new data
// @Tags			Owner
// @Accept			json
// @Produce		json
// @Param			ownerId		body		int		true	"OwnerId"
// @Param			name		body		string	false	"Name"
// @Param			surname		body		string	false	"Surname"
// @Param			patronymic	body		string	false	"Patronymic"
// @Success		200			{object}	Response
// @Failure		400			{object}	response.Response
// @Router			/owner/update [put]
func New(log *slog.Logger, ownerUpdater OwnerUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.owner.update.New"

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

		if req.Name != nil {
			err = ownerUpdater.UpdateOwnerName(req.OwnerId, *req.Name)
			if err != nil {
				log.Error("failed to update owner name", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update owner name"))

				return
			}
		}

		if req.Surname != nil {
			err = ownerUpdater.UpdateOwnerSurname(req.OwnerId, *req.Surname)
			if err != nil {
				log.Error("failed to update owner surname", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update owner surname"))

				return
			}
		}

		if req.Patronymic != nil {
			err = ownerUpdater.UpdateOwnerPatronymic(req.OwnerId, *req.Patronymic)
			if err != nil {
				log.Error("failed to update owner patronymic", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update owner patronymic"))

				return
			}
		}

		render.JSON(w, r, Response{
			response.OK(),
		})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if req.OwnerId < 1 {
		return false, slog.String("field", "owner_id"), "field owner_id is not valid"
	}
	if req.Name != nil && len(*req.Name) < 1 {
		return false, slog.String("field", "name"), "field name is not valid"
	}
	if req.Surname != nil && len(*req.Surname) < 1 {
		return false, slog.String("field", "surname"), "field surname is not valid"
	}
	if req.Patronymic != nil && len(*req.Patronymic) < 1 {
		return false, slog.String("field", "patronymic"), "field patronymic is not valid"
	}
	return true, slog.Attr{}, ""
}
