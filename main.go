package main

import (
	"escoutdoor/numguesser/logger"
	"log/slog"
)

func main() {
	logger.SetupLogger()
	server := NewServer(":8080")

	slog.Info("server is running", "port", server.listenAddress)
	if err := server.Start(); err != nil {
		panic(err)
	}
}
