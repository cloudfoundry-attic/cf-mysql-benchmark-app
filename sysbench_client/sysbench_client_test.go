package sysbench_client_test

import (
	"errors"
	"fmt"

	"database/sql"
	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client"
	fakeOsClient "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/os_client/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var _ = Describe("SysbenchClient", func() {
	const nodeCount = 3

	var (
		osClient       *fakeOsClient.FakeOsClient
		sysbenchClient sysbench_client.SysbenchClient
		config         conf.Config
		cmdName        string
		cmdArgs        []string
		nodeIndex      int
		cmdAction      string
		dbs            []*sql.DB
		mockDbs        []sqlmock.Sqlmock
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
			MySqlPort:        9999,
		}

		dbs = []*sql.DB{}
		mockDbs = []sqlmock.Sqlmock{}
		nodeIndex = 0
		for i := 0; i < nodeCount; i++ {
			db, mockDb, err := sqlmock.New()
			Expect(err).ToNot(HaveOccurred())
			dbs = append(dbs, db)
			mockDbs = append(mockDbs, mockDb)
		}

		sysbenchClient = sysbench_client.New(osClient, config, dbs)
	})

	JustBeforeEach(func() {
		cmdName = "sysbench"
		cmdArgs = []string{
			fmt.Sprintf("--mysql-host=%s", config.MySqlHosts[nodeIndex].Address),
			fmt.Sprintf("--mysql-port=%d", config.MySqlPort),
			fmt.Sprintf("--mysql-user=%s", config.MySqlUser),
			fmt.Sprintf("--mysql-password=%s", config.MySqlPwd),
			fmt.Sprintf("--mysql-db=%s", config.BenchmarkDB),
			fmt.Sprintf("--test=%s", "oltp"),
			fmt.Sprintf("--oltp-table-size=%d", config.NumBenchmarkRows),
			cmdAction,
		}
	})

	AfterEach(func() {
		for _, db := range dbs {
			db.Close()
		}
	})

	Describe("start", func() {
		BeforeEach(func() {
			cmdAction = "run"
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

	Describe("prepare", func() {

		BeforeEach(func() {
			cmdAction = "prepare"

			osClient.CombinedOutputReturns([]byte("Successfully ran"), nil)
		})

		Context("test table has the same number of rows set in the config", func() {
			It("does not run sysbench prepare", func() {
				mock := mockDbs[nodeIndex]

				mock.ExpectExec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", config.BenchmarkDB)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				tableRows := sqlmock.NewRows([]string{"name"}).AddRow("sbtest")
				mock.ExpectQuery(fmt.Sprintf("SHOW TABLES IN `%s` LIKE 'sbtest'", config.BenchmarkDB)).
					WillReturnRows(tableRows)

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(config.NumBenchmarkRows)
				// sqlmock interprets expects as a regex
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM .*`).WillReturnRows(countRows)

				_, err := sysbenchClient.Prepare(nodeIndex)
				Expect(err).ToNot(HaveOccurred())

				Expect(mock.ExpectationsWereMet()).To(Succeed())

				Expect(osClient.CombinedOutputCallCount()).To(Equal(0))
			})
		})

		Context("test table a different number of rows than the config", func() {
			It("truncates the table and runs sysbench prepare", func() {
				mock := mockDbs[nodeIndex]

				mock.ExpectExec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", config.BenchmarkDB)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				tableRows := sqlmock.NewRows([]string{"name"}).AddRow("sbtest")
				mock.ExpectQuery(fmt.Sprintf("SHOW TABLES IN `%s` LIKE 'sbtest'", config.BenchmarkDB)).
					WillReturnRows(tableRows)

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				// sqlmock interprets expects as a regex
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM .*`).WillReturnRows(countRows)

				mock.ExpectExec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`.sbtest", config.BenchmarkDB)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				_, err := sysbenchClient.Prepare(nodeIndex)
				Expect(err).ToNot(HaveOccurred())

				Expect(mock.ExpectationsWereMet()).To(Succeed())

				Expect(osClient.CombinedOutputCallCount()).To(Equal(1))
				actualName, actualArgs := osClient.CombinedOutputArgsForCall(0)
				Expect(actualName).To(Equal(cmdName))
				Expect(actualArgs).To(Equal(cmdArgs))
			})
		})

	})
})
