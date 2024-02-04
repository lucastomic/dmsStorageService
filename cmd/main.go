package main

import (
	"github.com/lucastomic/dmsStorageService/internal/controller"
	"github.com/lucastomic/dmsStorageService/internal/logging"
	"github.com/lucastomic/dmsStorageService/internal/middleware"
	"github.com/lucastomic/dmsStorageService/internal/pathrepository"
	"github.com/lucastomic/dmsStorageService/internal/pathservice"
	"github.com/lucastomic/dmsStorageService/internal/server"
	"github.com/lucastomic/dmsStorageService/internal/storageservice"
)

func main() {
	apilogger := logging.NewLogrusLogger()
	logicLogger := logging.NewLogrusLogger()
	dataLogger := logging.NewLogrusLogger()
	pathRepo := pathrepository.MemoryRepository(dataLogger)
	pathservice := pathservice.New(logicLogger, pathRepo)
	storageservice := storageservice.New(logicLogger, pathservice)
	controller := controller.New(logicLogger, storageservice)
	middlewares := []middleware.Middleware{
		middleware.NewLoggingMiddleware(apilogger),
		middleware.NewRequestIDMiddleware(),
		middleware.NewPathVarsMiddleware(),
	}
	server := server.New(":3003", controller, apilogger, logicLogger, middlewares)
	server.Run()
}
