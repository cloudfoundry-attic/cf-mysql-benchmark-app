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
			config = conf.NewConfig()
			err := config.ParseEnv()
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			os.Setenv("PROXY_IPS", "some-proxy-ip,another-proxy-ip")
			os.Setenv("BACKEND_IPS", "some-backend-ip,another-backend-ip")
			os.Setenv("ELB_IP", "some-elb-ip")
			os.Setenv("DATADOG_KEY", "some-datadog-key")
			os.Setenv("MYSQL_USER", "some-mysql-user")
			os.Setenv("MYSQL_PASSWORD", "some-mysql-password")
			os.Setenv("NUMBER_TEST_ROWS", "25")
			os.Setenv("TEST_DB", "some-db-name")
			os.Setenv("PORT", "9999")
		})

		AfterEach(func() {
			os.Unsetenv("PROXY_IPS")
			os.Unsetenv("BACKEND_IPS")
			os.Unsetenv("ELB_IP")
			os.Unsetenv("DATADOG_KEY")
			os.Unsetenv("MYSQL_USER")
			os.Unsetenv("MYSQL_PASSWORD")
			os.Unsetenv("NUMBER_TEST_ROWS")
			os.Unsetenv("TEST_DB")
			os.Unsetenv("PORT")
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
			err := test_helpers.IsRequiredField(config, "ProxyIPs")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if the Backend IPs array is empty", func() {
			err := test_helpers.IsRequiredField(config, "BackendIPs")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if Port is blank", func() {
			err := test_helpers.IsRequiredField(config, "Port")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
