package initializers

import (
	"log/slog"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "ravito"

func LoadEnvVars() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		slog.Error("Error loading .env file: " + err.Error())
	}
}
