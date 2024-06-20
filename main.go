package main

import "log/slog"

func main() {
	server := NewServer(":8080")

	slog.Info("server is running", "port", server.listenAddress)
	if err := server.Start(); err != nil {
		panic(err)
	}
}
