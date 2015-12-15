package config_test

import (
	"fmt"

	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/service-config/test_helpers"
)

var _ = Describe("Config", func() {
	Describe("Validate", func() {

		var (
			config    *conf.Config
			rawConfig string
		)

		JustBeforeEach(func() {
			osArgs := []string{
				"benchmarkApp",
				fmt.Sprintf("-config=%s", rawConfig),
			}

			var err error
			config, err = conf.NewConfig(osArgs)
			Expect(err).ToNot(HaveOccurred())
		})

		BeforeEach(func() {
			rawConfig = `{
				"ProxyIPs": ["some-proxy-ip", "another-proxy-ip"],
				"BackendIPs": ["some-backend-ip", "another-backend-ip"],
				"ElbIP": "some-elb-ip",
				"DatadogKey": "some-datadog",
				"MySqlUser": "some-username",
				"MySqlPwd": "some-password",
				"NumBenchmarkRows": 7,
				"BenchmarkDB": "some-db-name"
			}`
		})

		It("does not return error on valid config", func() {
			err := config.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error if the ELB IP is blank", func() {
			err := test_helpers.IsRequiredField(config, "ElbIP")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if MySQL user is blank", func() {
			err := test_helpers.IsRequiredField(config, "MySqlUser")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if MySQL password is blank", func() {
			err := test_helpers.IsRequiredField(config, "MySqlPwd")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if Benchmark DB name is blank", func() {
			err := test_helpers.IsRequiredField(config, "BenchmarkDB")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if the Proxy IPs array is empty", func() {
			config.ProxyIPs = []string{}
			err := config.Validate()
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the Backend IPs array is empty", func() {
			config.BackendIPs = []string{}
			err := config.Validate()
			Expect(err).To(HaveOccurred())
		})
	})
})
