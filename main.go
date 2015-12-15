package main

import (
	"net/http"

	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/api"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/os_client"
"fmt"
)

func main() {
	rootConfig := config.NewConfig()
	logger := rootConfig.Logger

	err := rootConfig.ParseEnv()
	if err != nil {
		logger.Fatal("Failed to parse environment variables", err)
	}

	err = rootConfig.Validate()
	if err != nil {
		logger.Fatal("Failed to validate config", err)
	}

	osClient := os_client.New()
	sysbenchClient := sysbench_client.New(osClient, *rootConfig)

	router, err := api.NewRouter(api.Api{
		RootConfig:     rootConfig,
		Routes:         api.DefaultRoutes(),
		SysbenchClient: sysbenchClient,
	})
	if err != nil {
		logger.Fatal("Failed to create router", err)
	}

	logger.Info("MySQL Benchmarking is running")
	err = http.ListenAndServe(fmt.Sprintf(":%d", rootConfig.Port), router)
	if err != nil {
		logger.Fatal("MySQL Benchmarking stopped unexpectedly", err)
	}
	logger.Info("MySQL Benchmarking has stopped gracefully")
}
