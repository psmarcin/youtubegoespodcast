package main

import (
	"ygp/pkg/api"
	"ygp/pkg/config"
	"ygp/pkg/db"
	"ygp/pkg/logger"
	"ygp/pkg/redis_client"
)

func main() {
	// Config
	config.Init()
	// Logger
	logger.Setup()
	// DB
	db.Setup()
	db.Migrate()
	defer db.Teardown()
	// Cache
	redis_client.Connect()
	defer redis_client.Teardown()
	// API
	api.Start()
}
