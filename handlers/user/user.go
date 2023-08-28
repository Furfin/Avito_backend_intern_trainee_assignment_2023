package user

import (
	"encoding/csv"
	"errors"
	"example/ravito/initializers"
	"example/ravito/models"
	"fmt"
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

type RequestUser struct {
	Userid int64 `json:"userid" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type UserResponse struct {
	Status   string `json:"status"`
	Error    string `json:"error"`
	Segments []models.Segment
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

			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", "Empty request"})

			return
		}
		if err != nil {
			log.Error("failed to decode request body")

			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", "Failed to decode request" + ": " + err.Error()})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		v := validator.New()

		if err := v.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request")

			render.Status(r, 400)
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
				render.Status(r, 400)
				render.JSON(w, r, Response{"Error", "Db creation problem"})

				return
			}
		}

		var seg models.Segment
		for _, val := range req.AddTo {
			seg = models.Segment{}
			if initializers.DB.Where(&seg, "slug = ?", val).Error != nil {
				log.Error("invalid request")
				render.Status(r, 400)
				render.JSON(w, r, Response{"Error", "Invalid segment slug: " + val})

				return
			}
		}
		for _, val := range req.RemoveFrom {
			seg = models.Segment{}
			if initializers.DB.Where(&seg, "slug = ?", val).Error != nil {
				log.Error("invalid request")
				render.Status(r, 400)
				render.JSON(w, r, Response{"Error", "Invalid segment slug: " + val})

				return
			}
		}
		var rel models.UserSegment
		for _, val := range req.AddTo {
			seg = models.Segment{}
			rel = models.UserSegment{}
			initializers.DB.First(&seg, "slug = ?", val)
			if initializers.DB.Where("user_id = ?", user.ID).First(&rel, "segment_id = ?", seg.ID).Error == nil {
				log.Error("invalid request")
				render.Status(r, 400)
				render.JSON(w, r, Response{"Error", strconv.FormatInt(user.Userid, 10) + " already in " + val})

				return
			}

			rel = models.UserSegment{UserID: int(user.ID), User: user, SegmentID: int(seg.ID), Segment: seg}
			result := initializers.DB.Create(&rel)
			if result.Error != nil {
				log.Error("invalid request")
				render.Status(r, 400)
				render.JSON(w, r, Response{"Error", "db creation problem"})

				return
			}
		}

		for _, val := range req.RemoveFrom {
			seg = models.Segment{}
			rel = models.UserSegment{}
			initializers.DB.First(&seg, "slug = ?", val)
			if initializers.DB.Where("user_id = ?", user.ID).First(&rel, "segment_id = ?", seg.ID).Error != nil {
				log.Error("invalid request")
				render.Status(r, 400)
				render.JSON(w, r, Response{"Error", strconv.FormatInt(user.Userid, 10) + " not in " + val})

				return
			}

			result := initializers.DB.Delete(&rel)
			if result.Error != nil {
				log.Error("invalid request")
				render.Status(r, 400)
				render.JSON(w, r, Response{"Error", "Db creation problem"})

				return
			}
		}

		render.JSON(w, r, Response{"Ok", "User segments changed"})
	}
}

func NewUser(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.add.NewUser"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestUser

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", "Empty request"})

			return
		}
		if err != nil {
			log.Error("failed to decode request body")
			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", "Failed to decode request" + ": " + err.Error()})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		v := validator.New()

		if err := v.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request")
			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", validateErr.Error()})

			return
		}

		// fmt.Println(req.AddTo[0])
		var user models.User
		var seg models.Segment
		var segs []models.Segment
		var rels []models.UserSegment

		if initializers.DB.First(&user, "userid = ?", req.Userid).Error != nil {
			log.Error("invalid request")
			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", "No such user is found"})
			return

		}

		initializers.DB.Where("user_id = ?", user.ID).Find(&rels)
		fmt.Println(rels, "ahhahahha")
		for _, rel := range rels {
			seg = models.Segment{}
			if initializers.DB.First(&seg, "id = ?", rel.SegmentID).Error == nil {
				segs = append(segs, seg)
			}
		}

		render.JSON(w, r, UserResponse{"Ok", "User info", segs})
	}
}

func NewUserHistory(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.add.NewUser"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestUser

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", "Empty request"})

			return
		}
		if err != nil {
			log.Error("failed to decode request body")
			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", "Failed to decode request" + ": " + err.Error()})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		v := validator.New()

		if err := v.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request")
			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", validateErr.Error()})

			return
		}

		// fmt.Println(req.AddTo[0])
		var user models.User
		var rels []models.UserSegment

		if initializers.DB.First(&user, "userid = ?", req.Userid).Error != nil {
			log.Error("invalid request")
			render.Status(r, 400)
			render.JSON(w, r, Response{"Error", "No such user is found"})
			return

		}
		records := [][]string{{"user", "segment_slug", "added/deleted", "datetime"}}
		initializers.DB.Unscoped().Where("user_id = ?", user.ID).Find(&rels)
		for _, rel := range rels {
			seg := models.Segment{}
			initializers.DB.Where("id = ?", rel.SegmentID).Find(&seg)
			records = append(records, []string{strconv.FormatInt(user.Userid, 10), seg.Slug, "added", rel.CreatedAt.String()})
		}
		rels = []models.UserSegment{}

		initializers.DB.Unscoped().Where("user_id = ?", user.ID).Where("deleted_at IS NOT NULL").Find(&rels)
		for _, rel := range rels {
			seg := models.Segment{}
			initializers.DB.Where("id = ?", rel.SegmentID).Find(&seg)
			records = append(records, []string{strconv.FormatInt(user.Userid, 10), seg.Slug, "deleted", rel.DeletedAt.Time.String()})
		}
		wr := csv.NewWriter(w)
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Add("Content-Disposition", `attachment; filename="history.csv"`)
		if err := wr.WriteAll(records); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
