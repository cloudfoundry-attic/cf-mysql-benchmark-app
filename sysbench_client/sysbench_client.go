package sysbench_client

import (
	"fmt"

	"database/sql"

	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/os_client"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pivotal-golang/lager"
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
	s.config.Logger.Info("Starting sysbench run", lager.Data{
		"MySQL Node": nodeIndex,
	})

	commandArgs := s.makeCommand(nodeIndex, "run")

	output, err := s.osClient.CombinedOutput("sysbench", commandArgs...)
	if err != nil {
		s.config.Logger.Error("Fatal - sysbench run failed", err, lager.Data{
			"MySQL Node": nodeIndex,
		})
		return string(output), fmt.Errorf("Sysbench failed to run! Error: %s", err.Error())
	}

	s.config.Logger.Info("Success - sysbench run succeeded", lager.Data{
		"MySQL Node": nodeIndex,
	})

	return string(output), nil
}

func (s sysbenchClient) Prepare(nodeIndex int) (string, error) {
	s.config.Logger.Info("Starting to prepare sysbench DB", lager.Data{
		"MySQL Node": nodeIndex,
	})

	db := s.dbs[nodeIndex]
	dbName := s.config.BenchmarkDB

	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		s.config.Logger.Error("Fatal - failed to create DB", err, lager.Data{
			"MySQL Node":    nodeIndex,
			"Database Name": dbName,
		})
		return "", fmt.Errorf("Database could not be created! Error: %s", err.Error())
	}

	dbIsTestReady, err := s.dbIsTestReady(db, nodeIndex)
	if err != nil {
		s.config.Logger.Error("Fatal - unable to determine state of DB", err, lager.Data{
			"MySQL Node":    nodeIndex,
			"Database Name": dbName,
		})
		return "", err
	}

	if dbIsTestReady == false {
		_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`.sbtest", dbName))
		if err != nil {
			return "", fmt.Errorf("Could not drop 'sbtest'! Error: %s", err.Error())
		}
		err = s.prepare(nodeIndex)
		if err != nil {
			s.config.Logger.Error("Fatal - ", err, lager.Data{
				"MySQL Node": nodeIndex,
			})
			return "", err
		}
		s.config.Logger.Info("Success - Sysbench prepared test DB ", lager.Data{
			"MySQL Node": nodeIndex,
		})
	}
	s.config.Logger.Info("Success - Node is ready for test", lager.Data{
		"MySQL Node": nodeIndex,
	})

	return fmt.Sprintf("Successfully prepared database: %s", dbName), nil
}

func (s sysbenchClient) dbIsTestReady(db *sql.DB, nodeIndex int) (bool, error) {

	dbName := s.config.BenchmarkDB

	var unused string
	err := db.QueryRow(fmt.Sprintf("SHOW TABLES IN `%s` LIKE 'sbtest'", dbName)).Scan(&unused)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	}

	// table does exist
	var rowCount int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s`.sbtest", dbName)).Scan(&rowCount)
	if err != nil {
		return false, err
	}

	rowCountMatches := (rowCount == s.config.NumBenchmarkRows)
	return rowCountMatches, nil
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
		fmt.Sprintf("--mysql-port=%d", s.config.MySqlPort),
		fmt.Sprintf("--mysql-user=%s", s.config.MySqlUser),
		fmt.Sprintf("--mysql-password=%s", s.config.MySqlPwd),
		fmt.Sprintf("--mysql-db=%s", s.config.BenchmarkDB),
		fmt.Sprintf("--test=%s", "oltp"),
		fmt.Sprintf("--oltp-table-size=%d", s.config.NumBenchmarkRows),
		fmt.Sprintf("--max-time=%d", s.config.MaxTime),
		fmt.Sprintf("--num-threads=%d", s.config.NumThreads),
	}
	return append(cmdArgs, sysbenchCommand)
}
