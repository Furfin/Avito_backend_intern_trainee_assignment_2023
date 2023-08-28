package main

import (
	"example/ravito/handlers/segment"
	"example/ravito/handlers/user"
	"example/ravito/initializers"
	"net/http"
	"os"
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

}
