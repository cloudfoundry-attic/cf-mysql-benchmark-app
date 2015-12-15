package sysbench_client

import (
	"fmt"

	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/os_client"
)

type SysbenchClient interface {
	Start(string) (string, error)
}

type sysbenchClient struct {
	osClient os_client.OsClient
	config   conf.Config
}

func New(osClient os_client.OsClient, config conf.Config) SysbenchClient {
	return &sysbenchClient{
		osClient: osClient,
		config:   config,
	}
}

func (s sysbenchClient) Start(nodeName string) (string, error) {
	commandArgs := s.makeCommand("run")

	err := s.osClient.Exec("sysbench", commandArgs...)
	if err != nil {
		return fmt.Sprintf("Sysbench failed to run! Error: %s", err.Error()), err
	}
	return fmt.Sprintf("Successfully ran test on node: %s", nodeName), nil
}

func (s sysbenchClient) makeCommand(sysbenchCommand string) []string {
	cmdArgs := []string{
		fmt.Sprintf("--mysql-host=%s", s.config.ElbIP),
		fmt.Sprintf("--mysql-port=%d", 3600),
		fmt.Sprintf("--mysql-user=%s", s.config.MySqlUser),
		fmt.Sprintf("--mysql-password=%s", s.config.MySqlPwd),
		fmt.Sprintf("--mysql-db=%s", s.config.BenchmarkDB),
		fmt.Sprintf("--test=%s", "oltp"),
		fmt.Sprintf("--oltp-table-size=%d", s.config.NumBenchmarkRows),
	}
	return append(cmdArgs, sysbenchCommand)
}
