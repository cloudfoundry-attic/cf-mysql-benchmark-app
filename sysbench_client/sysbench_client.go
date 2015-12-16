package sysbench_client

import (
	"fmt"

	"database/sql"
	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/os_client"
	_ "github.com/go-sql-driver/mysql"
)

type SysbenchClient interface {
	Start(int) (string, error)
	Prepare(int) (string, error)
}

type sysbenchClient struct {
	osClient os_client.OsClient
	config   conf.Config
	dbs      []*sql.DB
}

func New(osClient os_client.OsClient, config conf.Config, dbs []*sql.DB) SysbenchClient {
	return &sysbenchClient{
		osClient: osClient,
		config:   config,
		dbs:      dbs,
	}
}

func (s sysbenchClient) Start(nodeIndex int) (string, error) {
	commandArgs := s.makeCommand(nodeIndex, "run")

	output, err := s.osClient.CombinedOutput("sysbench", commandArgs...)
	if err != nil {
		return string(output), fmt.Errorf("Sysbench failed to run! Error: %s", err.Error())
	}
	return string(output), nil
}

func (s sysbenchClient) Prepare(nodeIndex int) (string, error) {
	db := s.dbs[nodeIndex]
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", s.config.BenchmarkDB))
	if err != nil {
		return "", fmt.Errorf("Database could not be created! Error: %s", err.Error())
	}

	// 'sbtest' is the default table name that sysbench creates for its testing.
	// There isn't any way to configure it differently, afawk.
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sbtest")
	if err != nil {
		return "", fmt.Errorf("Table 'sbtest' could not be created! Error: %s", err.Error())
	}

	var rowCount int
	err = db.QueryRow("SELECT COUNT(*) FROM sbtest").Scan(&rowCount)
	if err != nil {
		return "", fmt.Errorf("Failed to determine row count. Error: %s", err.Error())
	}

	if rowCount != s.config.NumBenchmarkRows {
		_, err = db.Exec("TRUNCATE TABLE sbtest")
		if err != nil {
			return "", fmt.Errorf("Could not truncate 'sbtest'! Error: %s", err.Error())
		}

		commandArgs := s.makeCommand(nodeIndex, "prepare")
		_, err = s.osClient.CombinedOutput("sysbench", commandArgs...)
		if err != nil {
			return "", fmt.Errorf("Sysbench failed to prepare! Error %s", err.Error())
		}
	}

	return "", nil
}

func (s sysbenchClient) makeCommand(nodeIndex int, sysbenchCommand string) []string {
	cmdArgs := []string{
		fmt.Sprintf("--mysql-host=%s", s.config.MySqlHosts[nodeIndex].Address),
		fmt.Sprintf("--mysql-port=%d", 3306),
		fmt.Sprintf("--mysql-user=%s", s.config.MySqlUser),
		fmt.Sprintf("--mysql-password=%s", s.config.MySqlPwd),
		fmt.Sprintf("--mysql-db=%s", s.config.BenchmarkDB),
		fmt.Sprintf("--test=%s", "oltp"),
		fmt.Sprintf("--oltp-table-size=%d", s.config.NumBenchmarkRows),
	}
	return append(cmdArgs, sysbenchCommand)
}
