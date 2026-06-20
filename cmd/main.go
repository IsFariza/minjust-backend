package main

import (
	"minjust-website/internal/app"
	"minjust-website/internal/config"
)

func main() {
	cfg := config.Load()
	app.Run(cfg)
}
