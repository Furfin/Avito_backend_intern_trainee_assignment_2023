package initializers

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func SetupLogger() *slog.Logger {

	Log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return Log
}
