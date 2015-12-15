package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/pivotal-golang/lager"
	"gopkg.in/validator.v2"
)

type Config struct {
	ProxyIPs         []string `validate:"min=1"`
	BackendIPs       []string `validate:"min=1"`
	ElbIP            string   `validate:"nonzero"`
	DatadogKey       string   `validate:"nonzero"`
	MySqlUser        string   `validate:"nonzero"`
	MySqlPwd         string   `validate:"nonzero"`
	NumBenchmarkRows int      `validate:"nonzero"`
	BenchmarkDB      string   `validate:"nonzero"`
	Logger           lager.Logger
}

func NewConfig() (*Config, error) {
	var rootConfig Config

	flags := flag.NewFlagSet("cf-mysql-benchmark", flag.ExitOnError)
	cf_lager.AddFlags(flags)
	rootConfig.Logger, _ = cf_lager.New("CF Mysql Benchmarking")

	benchmarkRows, err := strconv.Atoi(os.Getenv("NUMBER_TEST_ROWS"))
	if err != nil {
		return nil, err
	}
	rootConfig.ProxyIPs = strings.Split(os.Getenv("PROXY_IPS"), ",")
	rootConfig.BackendIPs = strings.Split(os.Getenv("BACKEND_IPS"), ",")
	rootConfig.ElbIP = os.Getenv("ELB_IP")
	rootConfig.DatadogKey = os.Getenv("DATADOG_KEY")
	rootConfig.MySqlUser = os.Getenv("BENCHMARK_MYSQL_USER")
	rootConfig.MySqlPwd = os.Getenv("BENCHMARK_MYSQL_PASSWORD")
	rootConfig.NumBenchmarkRows = benchmarkRows
	rootConfig.BenchmarkDB = os.Getenv("BENCHMARK_TEST_DB")

	return &rootConfig, nil
}

func (c Config) Validate() error {
	rootConfigErr := validator.Validate(c)
	var errString string
	if rootConfigErr != nil {
		errString = formatErrorString(rootConfigErr, "")
	}

	if len(errString) > 0 {
		return errors.New(fmt.Sprintf("Validation errors: %s\n", errString))
	}
	return nil
}

func formatErrorString(err error, keyPrefix string) string {
	errs := err.(validator.ErrorMap)
	var errsString string
	for fieldName, validationMessage := range errs {
		errsString += fmt.Sprintf("%s%s : %s\n", keyPrefix, fieldName, validationMessage)
	}
	return errsString
}
