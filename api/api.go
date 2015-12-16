package api

import (
	"net/http"

	"fmt"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client"
	"github.com/tedsuo/rata"
	"strconv"
)

type RunFunc func(nodeIndex int) (string, error)

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
		{Name: "prepare_test", Method: "POST", Path: "/prepare/:node"},
	}
}

func NewRouter(api Api) (http.Handler, error) {
	r := router{api: api}

	sysbenchClient := r.api.SysbenchClient
	api.Routes = DefaultRoutes()
	handlers := rata.Handlers{
		"start_test":   r.getInsecureHandler(sysbenchClient.Start),
		"prepare_test": r.getInsecureHandler(sysbenchClient.Prepare),
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
		targetNode, err := strconv.Atoi(rata.Param(req, "node"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMsg := fmt.Sprintf("Could not parse node index: %s", err.Error())
			w.Write([]byte(errMsg))
			return
		}

		body, err := run(targetNode)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.Write([]byte(body))
	})
}
