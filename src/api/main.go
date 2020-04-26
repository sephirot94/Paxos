package main

import(
	"github.com/gin-gonic/gin"
	"paxos/src/api/application"
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
