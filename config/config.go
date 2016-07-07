package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	//"strings"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/pivotal-golang/lager"
	"gopkg.in/validator.v2"
)

const DefaultMaxTime = 60
const DefaultNumThreads = 1
const DefaultBenchmarkRows = 100000

type Config struct {
	MySqlHost        MySqlHost   `validate`
	MySqlUser        string      `validate:"nonzero"`
	MySqlPassword    string      `validate:"nonzero"`
	MySqlPort        int         `validate:"nonzero"`
	DBName           string      `validate:"nonzero"`
	APIPort          int         `validate:"nonzero"`
	NumBenchmarkRows int         `validate:"nonzero"`
	MaxTime          int         `validate:"nonzero"`
	NumThreads       int         `validate:"nonzero"`
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

func (c *Config) ParseEnv() {
	c.MaxTime = stringToIntOrDefault(os.Getenv("MAX_TIME"), DefaultMaxTime)

	c.NumThreads, _ = strconv.Atoi(os.Getenv("NUM_THREADS"))
	if c.NumThreads == 0 {
		c.NumThreads = DefaultNumThreads
	}

	c.NumBenchmarkRows, _ = strconv.Atoi(os.Getenv("NUMBER_TEST_ROWS"))
	if c.NumBenchmarkRows == 0 {
		c.NumBenchmarkRows = DefaultBenchmarkRows
	}
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

func stringToIntOrDefault(stringValue string, otherwise int) int {
	result, err := strconv.Atoi(stringValue)

	if err != nil {
		result = otherwise
	}

	return result
}
