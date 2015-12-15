package sysbench_client_test

import (
	"errors"
	"fmt"

	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client"
	fakeOsClient "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/os_client/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SysbenchClient", func() {
	var (
		osClient       *fakeOsClient.FakeOsClient
		sysbenchClient sysbench_client.SysbenchClient
		config         conf.Config
		nodeName       string
		osCmd          []string
	)

	BeforeEach(func() {
		osClient = &fakeOsClient.FakeOsClient{}
		config = conf.Config{
			ElbIP:            "fake-elb-ip",
			MySqlUser:        "fake-mysql-user",
			MySqlPwd:         "fake-mysql-pwd",
			NumBenchmarkRows: 10,
			BenchmarkDB:      "fake-db",
		}
		sysbenchClient = sysbench_client.New(osClient, config)
		nodeName = "some-node"
	})

	Context("start", func() {
		BeforeEach(func() {
			osCmd = []string{
				"sysbench",
				fmt.Sprintf("--mysql-host=%s", config.ElbIP),
				fmt.Sprintf("--mysql-port=%d", 3600),
				fmt.Sprintf("--mysql-user=%s", config.MySqlUser),
				fmt.Sprintf("--mysql-password=%s", config.MySqlPwd),
				fmt.Sprintf("--mysql-db=%s", config.BenchmarkDB),
				fmt.Sprintf("--test=%s", "oltp"),
				fmt.Sprintf("--oltp-table-size=%d", config.NumBenchmarkRows),
				"run",
			}
		})

		Context("when sysbench exits 0", func() {
			BeforeEach(func() {
				osClient.ExecReturns(nil)
			})

			It("shells outs to the OS", func() {
				str, err := sysbenchClient.Start(nodeName)
				Expect(err).NotTo(HaveOccurred())
				Expect(str).To(ContainSubstring("Successfully ran"))
				Expect(osClient.ExecCallCount()).To(Equal(1))
				Expect(osClient.ExecArgsForCall(0)).To(Equal(osCmd))
			})
		})

		Context("when sysbench exits 1", func() {
			BeforeEach(func() {
				osClient.ExecReturns(errors.New("here's an error!"))
			})

			It("bubbles the error back up", func() {
				_, err := sysbenchClient.Start(nodeName)
				Expect(osClient.ExecArgsForCall(0)).To(Equal(osCmd))

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("here's an error!"))
			})
		})
	})
})
