package main

import (
	"example/ravito/handlers/segment"
	"example/ravito/handlers/user"
	"example/ravito/initializers"
	"example/ravito/models"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "example/ravito/docs"

	"example/ravito/httpSwaggerfix"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

//	@title			ravito
//	@version		1.0
//	@description	This is simple user segmentation service
//	@BasePath /
//	@host localhost:8084

func main() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
	initializers.SyncDB()
	initializers.SetupLogger()
	log := initializers.Log
	log.Info("Started ravito api")
	log.Debug("Debug enabled")

	ticker := time.NewTicker(15 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var rels []models.UserSegment
				initializers.DB.Find(&rels)
				for _, rel := range rels {
					if rel.DaysExpire != 0 {
						if int(time.Now().Sub(rel.CreatedAt)/24.0) < rel.DaysExpire {
							log.Info("Delete: " + strconv.FormatInt(int64(rel.ID), 10))
							initializers.DB.Delete(&rel)
						}
					}
				}
				log.Info("done")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/segment", segment.CreateSegment)
	router.Delete("/segment", segment.DeleteSegment)
	router.Post("/user/{userid}/add", user.UserSegmentsUpdate)
	router.Get("/user/{userid}", user.GetUserInfo)
	router.Post("/user/{userid}/csv", user.GetUserHistory)
	router.Get("/swagger/*", httpSwaggerfix.Handler(httpSwaggerfix.URL("doc.json")))

	log.Info("Swagger docs are running on http://" + os.Getenv("ADDRS") + ":" + os.Getenv("PORT") + "/swagger/index.html")

	log.Info("server started on " + os.Getenv("ADDRS") + ":" + os.Getenv("PORT"))

	server := &http.Server{
		Addr:         os.Getenv("ADDRS") + ":" + os.Getenv("PORT"),
		Handler:      router,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("stopping server")
	close(quit)

}
