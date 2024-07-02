package main

import (
	"driver_backend/internal/config"
)

func main() {
	// Init config: cleanenv
	cfg := config.MustLoad()

	// TODO: Init logger: slog

	// TODO: Init storage: sqlite

	// TODO: Init router: chi, "chi render"

	// TODO: Run server
}
