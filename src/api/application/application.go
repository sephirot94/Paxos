package application

import (
	"paxos/src/api/controllers"

)

type Application struct {
	RestHandler      *controllers.RestHandler

}

func Build() *Application {

	restHandler := controllers.NewRestHandler()

	return &Application{
		RestHandler:      restHandler,

	}
}
