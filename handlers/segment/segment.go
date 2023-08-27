package segment

import (
	"errors"
	"example/ravito/initializers"
	"example/ravito/models"
	"io"
	"net/http"

	"github.com/go-chi/render"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Slug string `json:"slug" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func NewCreate(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.create.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, Response{"Error", "empty request"})

			return
		}
		if err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, Response{"Error", "failed to decode request"})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		v := validator.New()

		if err := v.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request")

			render.JSON(w, r, Response{"Error", validateErr.Error()})

			return
		}

		seg := models.Segment{Slug: req.Slug}
		result := initializers.DB.Create(&seg)

		if result.Error != nil {
			log.Error("invalid request")

			render.JSON(w, r, Response{"Error", "db creation problem"})

			return
		}

		log.Info("segment created", req.Slug)

		render.JSON(w, r, Response{"Ok", "segment created"})
	}
}

func NewDelete(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.create.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, Response{"Error", "empty request"})

			return
		}
		if err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, Response{"Error", "failed to decode request"})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		v := validator.New()

		if err := v.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request")

			render.JSON(w, r, Response{"Error", validateErr.Error()})

			return
		}

		var seg models.Segment
		initializers.DB.First(&seg, "slug = ?", req.Slug)
		if initializers.DB.First(&seg, "slug = ?", req.Slug).Error != nil {
			log.Error("invalid request")

			render.JSON(w, r, Response{"Error", "Could not find segment to delete"})
			return
		}
		initializers.DB.Delete(&seg)

		log.Info("segment deleted", req.Slug)

		render.JSON(w, r, Response{"Ok", "segment deleted"})
	}
}
