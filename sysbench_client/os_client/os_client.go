package os_client

import "os/exec"

type OsClient interface {
	Exec(...string) error
}

type osClient struct{}

func New() *osClient {
	return &osClient{}
}

func (os osClient) Exec(cmd string, args ...string) error {
	newCmd := exec.Command(cmd, args...)
	err := newCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
