package os_client

import "os/exec"

type OsClient interface {
	CombinedOutput(string, ...string) ([]byte, error)
}

type osClient struct{}

func New() *osClient {
	return &osClient{}
}

func (os osClient) CombinedOutput(cmd string, args ...string) ([]byte, error) {
	newCmd := exec.Command(cmd, args...)
	return newCmd.CombinedOutput()
}
