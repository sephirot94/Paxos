package application

import (
	"Paxos/src/api/controllers/middlewares"
	"net/http"
	"github.com/gin-gonic/gin"
)

func MapURLs(app *Application, Router *gin.Engine) {

	Router.Use(middlewares.CORSMiddleware())


	// Add health check
	Router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})


	// NOTE: MIDDLEWARES CAN BE USED HERE FOR AUTHENTICATION, AUDIT, OR ANY OTHER NEED

	// GET
	Router.GET("/account", app.RestHandler.GetAccountBalance)

	group := Router.Group("/transactions")
	{
		group.GET("", app.RestHandler.GetTransactionHistory)
		group.GET("/:id", app.RestHandler.GetTransaction)
	}

	// POST
	Router.POST("/transactions", app.RestHandler.ExecTransaction)
}