package delete

import (
	"effective_mobile_test/internal/lib/api/response"
	"effective_mobile_test/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	OwnerId int `json:"ownerId"`
}

type Response struct {
	response.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=OwnerDeleter
type OwnerDeleter interface {
	DeleteOwner(OwnerId int) error
}

//	@Summary		Delete owner
//	@Description	Delete owner by ownerId
//	@Tags			Owner
//	@Accept			json
//	@Produce		json
//	@Param			ownerId	body		int	true	"OwnerId"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	response.Response
//	@Router			/owner/delete [delete]
func New(log *slog.Logger, ownerDeleter OwnerDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.owner.delete.New"

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

		err = ownerDeleter.DeleteOwner(req.OwnerId)
		if err != nil {
			log.Error("failed to delete owner", sl.Err(err))

			render.JSON(w, r, response.Error("failed to delete owner"))

			return
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
	return true, slog.Attr{}, ""
}
