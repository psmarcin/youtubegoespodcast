package main

import (
	"ytg/pkg/api"
	"ytg/pkg/config"
	"ytg/pkg/db"
	"ytg/pkg/logger"
	"ytg/pkg/redis_client"
)

func main() {
	config.Init()
	logger.Setup()
	db.Setup()
	redis_client.Connect()
	api.Start()
}
