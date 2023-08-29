package segment

import (
	"errors"
	"example/ravito/initializers"
	"example/ravito/models"
	"io"
	"math/rand"
	"net/http"

	"github.com/go-chi/render"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Slug  string `json:"slug" validate:"required"`
	UPadd int    `json:"upadd,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// CreateSegment - Creates segment using unique slug from body
// @description Creates new segment, upadd parameter sets number of percents of user send to the new segment
// @Tags ravito
// @Accept  json
// @Produce  json
// @Param request body segment.Request true "query params"
// @Success 200 {object} segment.Response "api response"
// @Router /segment [post]
func CreateSegment(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.segment.create.New"
	log := initializers.Log
	log = log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req Request

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
		render.JSON(w, r, Response{"Error", "Failed to decode request"})

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

	seg := models.Segment{Slug: req.Slug}
	result := initializers.DB.Create(&seg)

	if result.Error != nil {
		log.Error("invalid request")
		render.Status(r, 400)
		render.JSON(w, r, Response{"Error", "Db creation problem"})

		return
	}

	log.Info("segment created", req.Slug)

	if req.UPadd != 0 {
		var users []models.User
		initializers.DB.Find(&users)
		n := len(users)
		unum := (n * req.UPadd) / 100

		arr := generateUniqueRandomNumbers(unum, n)

		for _, val := range arr {
			rel := models.UserSegment{UserID: int(users[val].ID), User: users[val], SegmentID: int(seg.ID), Segment: seg}
			result := initializers.DB.Create(&rel)
			if result.Error != nil {
				log.Error("invalid request")
				render.Status(r, 400)
				render.JSON(w, r, Response{"Error", "db creation problem"})
				return
			}
		}
		log.Info("Users added")
	}

	render.JSON(w, r, Response{"Ok", "segment created"})
}

// DeleteSegment - Delete segment using unique slug from body
// @description Deletes segments and deletes all relations to this segment
// @Tags ravito
// @Accept  json
// @Produce  json
// @Param request body segment.Request true "query params"
// @Success 200 {object} Response "api response"
// @Router /segment [delete]
func DeleteSegment(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.segment.create.New"
	log := initializers.Log
	log = log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req Request

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
		render.JSON(w, r, Response{"Error", "Failed to decode request"})

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

	var seg models.Segment
	initializers.DB.First(&seg, "slug = ?", req.Slug)
	if initializers.DB.First(&seg, "slug = ?", req.Slug).Error != nil {
		log.Error("invalid request")
		render.Status(r, 400)
		render.JSON(w, r, Response{"Error", "Could not find segment to delete"})

		return
	}

	var rels []models.UserSegment

	initializers.DB.Where("segment_id = ?", seg.ID).Find(&rels)

	initializers.DB.Delete(&rels)

	initializers.DB.Unscoped().Delete(&seg)

	log.Info("segment deleted", req.Slug)

	render.JSON(w, r, Response{"Ok", "Segment deleted"})
}

func generateUniqueRandomNumbers(n, max int) []int {
	set := make(map[int]bool)
	var result []int
	for len(set) < n {
		value := rand.Intn(max)
		if !set[value] {
			set[value] = true
			result = append(result, value)
		}
	}
	return result
}
