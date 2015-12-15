package api

import (
	"net/http"

	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client"
	"github.com/tedsuo/rata"
)

type RunFunc func(node string) (string, error)

type Api struct {
	Routes         rata.Routes
	RootConfig     *config.Config
	SysbenchClient sysbench_client.SysbenchClient
}

type router struct {
	api Api
}

func DefaultRoutes() rata.Routes {
	return rata.Routes{
		{Name: "start_test", Method: "POST", Path: "/start/:node"},
	}
}

func NewRouter(api Api) (http.Handler, error) {
	r := router{api: api}

	sysbenchClient := r.api.SysbenchClient
	api.Routes = DefaultRoutes()
	handlers := rata.Handlers{
		"start_test": r.getInsecureHandler(sysbenchClient.Start),
	}

	handler, err := rata.NewRouter(api.Routes, handlers)
	if err != nil {
		api.RootConfig.Logger.Error("Error initializing router", err)
		return nil, err
	}

	return handler, nil
}

func (r router) getInsecureHandler(run RunFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		targetNode := rata.Param(req, "node")
		body, err := run(targetNode)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(body))
	})
}
