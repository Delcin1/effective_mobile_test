package search

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
	postgres.SearchRequest
}

type Response struct {
	response.Response
	Cars []postgres.Car `json:"cars"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=CarSearcher
type CarSearcher interface {
	GetCarsBySearchRequest(searchRequest postgres.SearchRequest) ([]postgres.Car, error)
}

//	@Summary		Search cars
//	@Description	Search cars by search request
//	@Tags			Car
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	response.Response
//	@Router			/car/search [get]
func New(log *slog.Logger, carSearcher CarSearcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.car.search.New"

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

		cars, err := carSearcher.GetCarsBySearchRequest(req.SearchRequest)
		if err != nil {
			log.Error("failed to get cars by search request", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get cars by search request"))

			return
		}

		render.JSON(w, r, Response{
			response.OK(),
			cars,
		})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if req.PageSize < 1 {
		return false, slog.String("field", "pageSize"), "field pageSize is not valid"
	}
	if req.PageNum < 1 {
		return false, slog.String("field", "pageNum"), "field pageNum is not valid"
	}
	return true, slog.Attr{}, ""
}
