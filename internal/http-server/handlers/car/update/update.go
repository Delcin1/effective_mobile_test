package update

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
	CarId  int             `json:"carId"`
	RegNum *string         `json:"regNum,omitempty"`
	Mark   *string         `json:"mark,omitempty"`
	Model  *string         `json:"model,omitempty"`
	Year   *int            `json:"year,omitempty"`
	Owner  *postgres.Owner `json:"owner,omitempty"`
}

type Response struct {
	response.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=CarUpdater
type CarUpdater interface {
	UpdateRegNum(carID int, newRegNum string) error
	UpdateMark(carID int, newMark string) error
	UpdateModel(carID int, newModel string) error
	UpdateYear(carID int, newYear int) error
	UpdateOwner(carID int, newOwner postgres.Owner) error
}

//	@Summary		Update car
//	@Description	Update car by carId and new data
//	@Tags			Car
//	@Accept			json
//	@Produce		json
//	@Param			regNum	body		string			true	"RegNum"
//	@Param			mark	body		string			false	"Mark"
//	@Param			model	body		string			false	"Model"
//	@Param			year	body		int				false	"Year"
//	@Param			owner	body		postgres.Owner	false	"Owner"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	response.Response
//	@Router			/car/update [put]
func New(log *slog.Logger, carUpdater CarUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.car.update.New"

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

		if req.RegNum != nil {
			err = carUpdater.UpdateRegNum(req.CarId, *req.RegNum)
			if err != nil {
				log.Error("failed to update regNum", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update regNum"))

				return
			}
		}

		if req.Mark != nil {
			err = carUpdater.UpdateMark(req.CarId, *req.Mark)
			if err != nil {
				log.Error("failed to update mark", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update mark"))

				return
			}
		}

		if req.Model != nil {
			err = carUpdater.UpdateModel(req.CarId, *req.Model)
			if err != nil {
				log.Error("failed to update model", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update model"))

				return
			}
		}

		if req.Year != nil {
			err = carUpdater.UpdateYear(req.CarId, *req.Year)
			if err != nil {
				log.Error("failed to update year", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update year"))

				return
			}
		}

		if req.Owner != nil {
			err = carUpdater.UpdateOwner(req.CarId, *req.Owner)
			if err != nil {
				log.Error("failed to update owner", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update owner"))

				return
			}
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
	if req.RegNum != nil && len(*req.RegNum) < 1 {
		return false, slog.String("field", "regNum"), "field regNum is not valid"
	}
	if req.Mark != nil && len(*req.Mark) < 1 {
		return false, slog.String("field", "mark"), "field mark is not valid"
	}
	if req.Model != nil && len(*req.Model) < 1 {
		return false, slog.String("field", "model"), "field model is not valid"
	}
	if req.Year != nil && *req.Year < 1900 {
		return false, slog.String("field", "year"), "field year is not valid"
	}
	if req.Owner != nil && len(req.Owner.Name) < 1 {
		return false, slog.String("field", "owner.name"), "field owner.name is not valid"
	}
	if req.Owner != nil && len(req.Owner.Surname) < 1 {
		return false, slog.String("field", "owner.surname"), "field owner.surname is not valid"
	}
	if req.Owner != nil && len(req.Owner.Patronymic) < 1 {
		return false, slog.String("field", "owner.patronymic"), "field owner.patronymic is not valid"
	}
	return true, slog.Attr{}, ""
}
