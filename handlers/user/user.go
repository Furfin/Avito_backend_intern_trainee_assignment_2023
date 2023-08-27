package user

import (
	"errors"
	"example/ravito/initializers"
	"example/ravito/models"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type RequestUserSlugs struct {
	AddTo      []string `json:"AddTo,omitempty"`
	RemoveFrom []string `json:"RemoveFrom,omitempty"`
	Userid     int64    `json:"userid" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.add.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestUserSlugs

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, Response{"Error", "empty request"})

			return
		}
		if err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, Response{"Error", "failed to decode request" + ": " + err.Error()})

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

		// fmt.Println(req.AddTo[0])
		var user models.User

		if initializers.DB.First(&user, "userid = ?", req.Userid).Error != nil {
			user = models.User{Userid: req.Userid}
			result := initializers.DB.Create(&user)
			if result.Error != nil {
				log.Error("invalid request")
				render.JSON(w, r, Response{"Error", "db creation problem"})
				return
			}
		}

		var seg models.Segment
		for _, val := range req.AddTo {
			if initializers.DB.First(&seg, "slug = ?", val).Error != nil {
				log.Error("invalid request")
				render.JSON(w, r, Response{"Error", "Invalid segment slug: " + val})
				return
			}
		}
		for _, val := range req.RemoveFrom {
			if initializers.DB.First(&seg, "slug = ?", val).Error != nil {
				log.Error("invalid request")
				render.JSON(w, r, Response{"Error", "Invalid segment slug: " + val})
				return
			}
		}
		var rel models.UserSegment
		for _, val := range req.AddTo {
			initializers.DB.First(&seg, "slug = ?", val)
			if initializers.DB.Where("user_id = ?", user.ID).First(&rel, "segment_id = ?", seg.ID).Error == nil {
				log.Error("invalid request")
				render.JSON(w, r, Response{"Error", strconv.FormatInt(user.Userid, 10) + " already in " + val})
				return
			}

			rel = models.UserSegment{UserID: int(user.ID), User: user, SegmentID: int(seg.ID), Segment: seg}
			result := initializers.DB.Create(&rel)
			if result.Error != nil {
				log.Error("invalid request")
				render.JSON(w, r, Response{"Error", "db creation problem"})
				return
			}
		}

		for _, val := range req.RemoveFrom {
			initializers.DB.First(&seg, "slug = ?", val)
			if initializers.DB.Where("user_id = ?", user.ID).First(&rel, "segment_id = ?", seg.ID).Error != nil {
				log.Error("invalid request")
				render.JSON(w, r, Response{"Error", strconv.FormatInt(user.Userid, 10) + " not in " + val})
				return
			}

			result := initializers.DB.Delete(&rel)
			if result.Error != nil {
				log.Error("invalid request")
				render.JSON(w, r, Response{"Error", "db creation problem"})
				return
			}
		}

		render.JSON(w, r, Response{"Ok", "User test"})
	}
}
