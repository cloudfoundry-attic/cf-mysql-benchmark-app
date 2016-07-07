package config_test

import (
	"os"

	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"github.com/pivotal-cf-experimental/service-config/test_helpers"
)

var _ = Describe("Config", func() {
	Describe("Validate", func() {

		var (
			config *conf.Config
		)

		JustBeforeEach(func() {
			config = conf.NewConfig()
			config.ParseEnv()
		})

		BeforeEach(func() {
			//os.Setenv("MYSQL_HOSTS", "backend1=1.1.1.1,proxy1=2.2.2.2,elb=some.dns.name")
			//os.Setenv("MYSQL_USER", "some-mysql-user")
			//os.Setenv("MYSQL_PASSWORD", "some-mysql-password")
			//os.Setenv("TEST_DB", "some-db-name")
			//os.Setenv("PORT", "9999")
		})

		AfterEach(func() {
			//os.Unsetenv("MYSQL_HOSTS")
			//os.Unsetenv("MYSQL_USER")
			//os.Unsetenv("MYSQL_PASSWORD")
			//os.Unsetenv("MYSQL_PORT")
			//os.Unsetenv("TEST_DB")
			//os.Unsetenv("PORT")
		})

		It("does not return error on valid config", func() {
			err := config.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("Test Max Time", func() {
			Context("when there is no environment variable", func() {
				BeforeEach(func() {
					os.Unsetenv("MAX_TIME")
				})
				It("defaults to 60", func() {
					Expect(config.MaxTime).To(Equal(60))
				})
			})

			Context("when there is an environment variable", func() {
				BeforeEach(func() {
					os.Setenv("MAX_TIME", "1234")
				})
				It("sets it to the specified value", func() {
					Expect(config.MaxTime).To(Equal(1234))
				})
			})
		})

		Context("Setting the oltp row count", func() {
 			Context("when there is no environment variable", func() {
				BeforeEach(func() {
					os.Unsetenv("NUMBER_TEST_ROWS")
				})
				It("defaults to 100000", func() {
					Expect(config.NumBenchmarkRows).To(Equal(100000))
				})
			})

			Context("when there is an environment variable", func() {
				BeforeEach(func() {
					os.Setenv("NUMBER_TEST_ROWS", "4321")
				})
				It("sets it to the specified value", func() {
					Expect(config.NumBenchmarkRows).To(Equal(4321))
				})
			})
		})

		Context("Number of threads", func() {
			Context("when there is no environment variable", func() {
				BeforeEach(func() {
					os.Unsetenv("NUM_THREADS")
				})
				It("defaults to 1", func() {
					Expect(config.NumThreads).To(Equal(1))
				})
			})

			Context("when there is an environment variable", func() {
				BeforeEach(func() {
					os.Setenv("NUM_THREADS", "1234")
				})
				It("sets it to the specified value", func() {
					Expect(config.NumThreads).To(Equal(1234))
				})
			})
		})
	})
})
