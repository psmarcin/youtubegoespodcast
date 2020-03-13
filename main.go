package main

import (
	"ygp/pkg/api"
	"ygp/pkg/cache"
	"ygp/pkg/config"
	"ygp/pkg/logger"
)

func main() {
	// Config
	config.Init()
	// Logger
	logger.Setup()
	// Cache
	_, _ = cache.Connect()
	// API
	api.Start()
}
