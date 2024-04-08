package save

import (
	client2 "effective_mobile_test/internal/client"
	"effective_mobile_test/internal/lib/api/response"
	"effective_mobile_test/internal/lib/logger/sl"
	"effective_mobile_test/internal/storage/postgres"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	RegNums []string `json:"regNums"`
}

type Response struct {
	response.Response
	CarsIds []int `json:"cars_ids"`
}

type CarSaver interface {
	SaveCar(car postgres.Car) (int, error)
}

//	@Summary		Save a new car
//	@Description	Save a new car by regNums
//	@Tags			Car
//	@Accept			json
//	@Produce		json
//	@Param			regNums	body		[]string	true	"RegNums"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	response.Response
//	@Router			/car/save [post]
func New(log *slog.Logger, carSaver CarSaver, helpAPIUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.car.save.New"

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

		client := client2.SearchClient{URL: helpAPIUrl}
		resp, err := client.FindUsers(client2.SearchRequest{RegNums: req.RegNums})
		if err != nil {
			log.Error("failed to find car", sl.Err(err))

			render.JSON(w, r, response.Error("failed to find car"))
		}

		var carsIds []int
		for _, car := range resp.Cars {
			carId, err := carSaver.SaveCar(car)
			if err != nil {
				log.Error("failed to save car", sl.Err(err))

				render.JSON(w, r, response.Error("failed to save car"))
			}

			log.Info("car saved", slog.Int("car_id", carId))

			carsIds = append(carsIds, carId)
		}

		render.JSON(w, r, Response{
			response.OK(),
			carsIds,
		})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	for _, regNum := range req.RegNums {
		if len(regNum) < 1 || len(regNum) > 255 {
			return false, slog.String("field", "regNums"), "field regNums is not valid"
		}
	}
	return true, slog.Attr{}, ""
}
