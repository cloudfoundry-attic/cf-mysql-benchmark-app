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
			os.Setenv("DATADOG_KEY", "some-datadog-key")
			os.Setenv("MYSQL_HOSTS", "backend1=1.1.1.1,proxy1=2.2.2.2,elb=some.dns.name")
			os.Setenv("MYSQL_USER", "some-mysql-user")
			os.Setenv("MYSQL_PASSWORD", "some-mysql-password")
			os.Setenv("NUMBER_TEST_ROWS", "25")
			os.Setenv("TEST_DB", "some-db-name")
			os.Setenv("PORT", "9999")
		})

		AfterEach(func() {
			os.Unsetenv("DATADOG_KEY")
			os.Unsetenv("MYSQL_HOSTS")
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

		It("returns an error if MySQL Hosts is blank", func() {
			err := test_helpers.IsRequiredField(config, "MySqlHosts")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if MySqlHosts.Name is blank", func() {
			err := test_helpers.IsRequiredField(config, "MySqlHosts.Name")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if MySqlHosts.Address is blank", func() {
			err := test_helpers.IsRequiredField(config, "MySqlHosts.Address")
			Expect(err).ToNot(HaveOccurred())
		})

		It("parses MySQL Hosts into key value pair", func() {
			Expect(config.MySqlHosts).To(ConsistOf([]conf.MySqlHost{
				conf.MySqlHost{
					Name:    "backend1",
					Address: "1.1.1.1",
				},
				conf.MySqlHost{
					Name:    "proxy1",
					Address: "2.2.2.2",
				},
				conf.MySqlHost{
					Name:    "elb",
					Address: "some.dns.name",
				},
			}))
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

		It("returns an error if Port is blank", func() {
			err := test_helpers.IsRequiredField(config, "Port")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
