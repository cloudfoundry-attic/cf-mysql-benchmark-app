package sysbench_client

import (
	"fmt"

	conf "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/config"
	"github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/sysbench_client/os_client"
)

type SysbenchClient interface {
	Start(int) (string, error)
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

func (s sysbenchClient) Start(nodeIndex int) (string, error) {
	commandArgs := s.makeCommand(nodeIndex, "run")

	output, err := s.osClient.CombinedOutput("sysbench", commandArgs...)
	if err != nil {
		return string(output), fmt.Errorf("Sysbench failed to run! Error: %s", err.Error())
	}
	return string(output), nil
}

func (s sysbenchClient) makeCommand(nodeIndex int, sysbenchCommand string) []string {
	cmdArgs := []string{
		fmt.Sprintf("--mysql-port=%d", 3600),
		fmt.Sprintf("--mysql-host=%s", s.config.MySqlHosts[nodeIndex].Address),
		fmt.Sprintf("--mysql-user=%s", s.config.MySqlUser),
		fmt.Sprintf("--mysql-password=%s", s.config.MySqlPwd),
		fmt.Sprintf("--mysql-db=%s", s.config.BenchmarkDB),
		fmt.Sprintf("--test=%s", "oltp"),
		fmt.Sprintf("--oltp-table-size=%d", s.config.NumBenchmarkRows),
	}
	return append(cmdArgs, sysbenchCommand)
}
