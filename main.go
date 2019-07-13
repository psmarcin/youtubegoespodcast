package main

import (
	"ytg/pkg/api"
	"ytg/pkg/config"
	"ytg/pkg/db"
	"ytg/pkg/logger"
	"ytg/pkg/redis_client"
)

func main() {
	// Config
	config.Init()
	// Logger
	logger.Setup()
	// DB
	db.Setup()
	defer db.Teardown()
	// Cache
	redis_client.Connect()
	defer redis_client.Teardown()
	// API
	api.Start()
}
