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
	CarId int `json:"carId"`
}

type Response struct {
	response.Response
}

type CarDeleter interface {
	DeleteCar(carID int) error
}

//	@Summary		Delete car
//	@Description	Delete car by carId
//	@Tags			Car
//	@Accept			json
//	@Produce		json
//	@Param			carId	body		int	true	"CarId"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	response.Response
//	@Router			/car/delete [delete]
func New(log *slog.Logger, carDeleter CarDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.car.delete.New"

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

		err = carDeleter.DeleteCar(req.CarId)
		if err != nil {
			log.Error("failed to delete car", sl.Err(err))

			render.JSON(w, r, response.Error("failed to delete car"))

			return
		}

		render.JSON(w, r, Response{
			response.OK(),
		})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if req.CarId < 1 {
		return false, slog.String("field", "car_id"), "field car_id is not valid"
	}
	return true, slog.Attr{}, ""
}
