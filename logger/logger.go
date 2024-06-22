package logger

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

func SetupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	switch strings.ToLower(os.Getenv("LEVEL")) {
	case "debug", "d":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "error", "e":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "info", "i":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}
