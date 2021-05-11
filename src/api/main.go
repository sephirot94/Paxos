package main

import (
	"Paxos/src/api/application"
	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func main() {
	configureRouter()
}

func configureRouter() {
	Router = gin.Default()

	app := application.Build()
	application.MapURLs(app, Router)

	Router.Run("localhost:8080")
}
