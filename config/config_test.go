package config_test

import (
	"fmt"
    "errors"
    . "github.com/cf-mysql-benchmark-app/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/service-config/test_helpers"
)

var _ = Describe("Config", func() {
	Describe("Validate", func() {

		var (
			config *Config
			rawConfig string
		)

		JustBeforeEach(func() {
			osArgs := []string{
				fmt.Sprintf("-config=%s", rawConfig),
			}

			var err error
			config, err = NewConfig(osArgs)
			Expect(err).ToNot(HaveOccurred())
		})

		BeforeEach(func() {
            rawConfig := `{
                "ProxyIPs": ["10.10.163.11", "10.10.164.11"],
				"BackendIPs": ["10.10.163.10", "10.10.164.10", "10.10.165.10"],
				"ElbIP": "internal-cf-mysql-benchmarking-lb-1249448519.us-east-1.elb.amazonaws.com",
				"MySqlUser": "root",
				"MySqlPwd": "dusky7dirge",
				"NumBenchmarkRows": 10,
				"BenchmarkDB": "sysbench_db",
			}`
		})

		It("does not return error on valid config", func() {
			err := config.Validate()
			Expect(err).NotTo(HaveOccurred())
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
            config.ProxyIPs = []string
            Expect(err).To(HaveOccurred())
        })

        It("returns an error if the Backend IPs array is empty", func() {
            config.BackendIPs = []string
            Expect(err).To(HaveOccurred())
        })

        It("returns an error if the ELB IP is blank", func() {
            err := test_helpers.IsRequiredField(config, "ElbIP")
            Expect(err).To(HaveOccurred())
        })
	})
})