package config_test

import (
	"os"

	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/service-config/test_helpers"
)

var _ = Describe("Config", func() {
	Describe("Validate", func() {

		var (
			config *conf.Config
		)

		JustBeforeEach(func() {
			var err error
			config, err = conf.NewConfig()
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			os.Setenv("PROXY_IPS", "some-proxy-ip,another-proxy-ip")
			os.Setenv("BACKEND_IPS", "some-backend-ip,another-backend-ip")
			os.Setenv("ELB_IP", "some-elb-ip")
			os.Setenv("DATADOG_KEY", "some-datadog-key")
			os.Setenv("BENCHMARK_MYSQL_USER", "some-mysql-user")
			os.Setenv("BENCHMARK_MYSQL_PASSWORD", "some-mysql-password")
			os.Setenv("NUMBER_TEST_ROWS", "25")
			os.Setenv("BENCHMARK_TEST_DB", "some-db-name")
		})

		AfterEach(func() {
			os.Setenv("PROXY_IPS", "")
			os.Setenv("BACKEND_IPS", "")
			os.Setenv("ELB_IP", "")
			os.Setenv("DATADOG_KEY", "")
			os.Setenv("BENCHMARK_MYSQL_USER", "")
			os.Setenv("BENCHMARK_MYSQL_PASSWORD", "")
			os.Setenv("NUMBER_TEST_ROWS", "")
			os.Setenv("BENCHMARK_TEST_DB", "")
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
