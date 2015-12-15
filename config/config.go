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
	MySqlHosts       []MySqlHost `validate:"min=1"`
	DatadogKey       string      `validate:"nonzero"`
	MySqlUser        string      `validate:"nonzero"`
	MySqlPwd         string      `validate:"nonzero"`
	NumBenchmarkRows int         `validate:"nonzero"`
	BenchmarkDB      string      `validate:"nonzero"`
	Port             int         `validate:"nonzero"`
	Logger           lager.Logger
}

type MySqlHost struct {
	Name    string `validate:"nonzero"`
	Address string `validate:"nonzero"`
}

func NewConfig() *Config {
	rootConfig := Config{}

	flags := flag.NewFlagSet("cf-mysql-benchmark", flag.ExitOnError)
	cf_lager.AddFlags(flags)
	rootConfig.Logger, _ = cf_lager.New("CF Mysql Benchmarking")

	return &rootConfig
}

func (c *Config) ParseEnv() error {

	// will default to 0 if strconv fails
	c.NumBenchmarkRows, _ = strconv.Atoi(os.Getenv("NUMBER_TEST_ROWS"))
	c.Port, _ = strconv.Atoi(os.Getenv("PORT"))

	c.MySqlHosts = []MySqlHost{}
	hosts := strings.Split(os.Getenv("MYSQL_HOSTS"), ",")
	for _, host := range hosts {
		newHost := MySqlHost{
			Name:    strings.Split(host, "=")[0],
			Address: strings.Split(host, "=")[1],
		}
		c.MySqlHosts = append(c.MySqlHosts, newHost)
	}
	c.DatadogKey = os.Getenv("DATADOG_KEY")
	c.MySqlUser = os.Getenv("MYSQL_USER")
	c.MySqlPwd = os.Getenv("MYSQL_PASSWORD")
	c.BenchmarkDB = os.Getenv("TEST_DB")

	return nil
}

func (c Config) Validate() error {
	rootConfigErr := validator.Validate(c)
	var errString string
	if rootConfigErr != nil {
		errString = formatErrorString(rootConfigErr, "")
	}

	// validator.Validate does not work on nested arrays
	for i, host := range c.MySqlHosts {
		nestedErr := validator.Validate(host)
		if nestedErr != nil {
			errString += formatErrorString(
				nestedErr,
				fmt.Sprintf("MySqlHosts[%d].", i),
			)
		}
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
