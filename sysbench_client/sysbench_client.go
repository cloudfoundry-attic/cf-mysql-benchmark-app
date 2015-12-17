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
	dbName := s.config.BenchmarkDB

	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		return "", fmt.Errorf("Database could not be created! Error: %s", err.Error())
	}

	var unused string
	err = db.QueryRow(fmt.Sprintf("SHOW TABLES IN `%s` LIKE 'sbtest'", dbName)).Scan(&unused)

	switch {
	case err == nil:
		err = s.truncateTable(db, nodeIndex)
		if err != nil {
			return "", err
		}
	case err == sql.ErrNoRows:
		err = s.prepare(nodeIndex)
		if err != nil {
			return "", err
		}
	case err != nil:
		return "", err
	}

	return "", nil
}

func (s sysbenchClient) truncateTable(db *sql.DB, nodeIndex int) error {
	dbName := s.config.BenchmarkDB

	var rowCount int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s`.sbtest", dbName)).Scan(&rowCount)
	if err != nil {
		return fmt.Errorf("Failed to determine row count. Error: %s", err.Error())
	}

	if rowCount != s.config.NumBenchmarkRows {
		_, err = db.Exec(fmt.Sprintf("DROP TABLE `%s`.sbtest", dbName))
		if err != nil {
			return fmt.Errorf("Could not drop 'sbtest'! Error: %s", err.Error())
		}
		err = s.prepare(nodeIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s sysbenchClient) prepare(nodeIndex int) error {
	commandArgs := s.makeCommand(nodeIndex, "prepare")
	output, err := s.osClient.CombinedOutput("sysbench", commandArgs...)
	if err != nil {
		return fmt.Errorf("Sysbench failed to prepare! Error %s, Output: %s", err.Error(), output)
	}
	return nil
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
