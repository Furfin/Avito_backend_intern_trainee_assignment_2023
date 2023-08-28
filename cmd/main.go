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

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
	initializers.SyncDB()

	log := initializers.SetupLogger()
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

	router.Post("/segment", segment.NewCreate(log))
	router.Delete("/segment", segment.NewDelete(log))
	router.Post("/user", user.New(log))
	router.Get("/user", user.NewUser(log))
	router.Get("/user/csv", user.NewUserHistory(log))

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
