package main

import (
	"ytg/pkg/api"
	"ytg/pkg/config"
	"ytg/pkg/logger"
)

func main() {
	config.Init()
	logger.Setup()
	api.Start()
}
