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
		cmdName        string
		cmdArgs        []string
		nodeIndex      int
	)

	BeforeEach(func() {
		osClient = &fakeOsClient.FakeOsClient{}
		config = conf.Config{
			MySqlHosts: []conf.MySqlHost{
				{
					Name:    "host-0",
					Address: "1.1.1.1",
				},
				{
					Name:    "host-1",
					Address: "2.2.2.2",
				},
				{
					Name:    "host-2",
					Address: "3.3.3.3",
				},
			},
			MySqlUser:        "fake-mysql-user",
			MySqlPwd:         "fake-mysql-pwd",
			NumBenchmarkRows: 10,
			BenchmarkDB:      "fake-db",
		}
		sysbenchClient = sysbench_client.New(osClient, config)
		nodeIndex = 0
	})

	Context("start", func() {
		BeforeEach(func() {
			cmdName = "sysbench"
			cmdArgs = []string{
				fmt.Sprintf("--mysql-host=%s", config.MySqlHosts[nodeIndex].Address),
				fmt.Sprintf("--mysql-port=%d", 3306),
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
				osClient.CombinedOutputReturns([]byte("Successfully ran"), nil)
			})

			It("shells outs to the OS", func() {
				cmdOutput, err := sysbenchClient.Start(nodeIndex)
				Expect(err).NotTo(HaveOccurred())
				Expect(cmdOutput).To(ContainSubstring("Successfully ran"))
				Expect(osClient.CombinedOutputCallCount()).To(Equal(1))
				actualName, actualArgs := osClient.CombinedOutputArgsForCall(0)
				Expect(actualName).To(Equal(cmdName))
				Expect(actualArgs).To(Equal(cmdArgs))
			})
		})

		Context("when sysbench exits 1", func() {
			BeforeEach(func() {
				osClient.CombinedOutputReturns([]byte("fake-stderr"), errors.New("here's an error!"))
			})

			It("bubbles the error back up", func() {
				output, err := sysbenchClient.Start(nodeIndex)
				actualName, actualArgs := osClient.CombinedOutputArgsForCall(0)
				Expect(actualName).To(Equal(cmdName))
				Expect(actualArgs).To(Equal(cmdArgs))

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("here's an error!"))
				Expect(output).To(ContainSubstring("fake-stderr"))
			})
		})
	})
})
