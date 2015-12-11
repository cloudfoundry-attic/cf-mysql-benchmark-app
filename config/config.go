package config

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/pivotal-cf-experimental/service-config"
	"github.com/pivotal-golang/lager"
	"gopkg.in/validator.v2"
)

type Config struct {
	ProxyIPs     []string
    BackendIPs   []string
    ElbIP        []string
    DatadogKey   string
    MySqlUser    string
    MySqlPwd     string
    NumBenchmarkRows int
    BenchmarkDB  string
    Logger       lager.Logger
}

func NewConfig(osArgs []string) (*Config, error) {
	var config Config

	configurationOptions := osArgs[0:]

	serviceConfig := service_config.New()
	flags := flag.NewFlagSet("benchmarkApp", flag.ExitOnError)

	cf_lager.AddFlags(flags)

	serviceConfig.AddFlags(flags)
	flags.Parse(configurationOptions)

	err := serviceConfig.Read(&config)

	config.Logger, _ = cf_lager.New("benchmarkApp")

    return &config, err
}

func (c Config) Validate() error {
	/*configErr := validator.Validate(c)
	var errString string
	if configErr != nil {
		errString = formatErrorString(configErr, "")
	}*/

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