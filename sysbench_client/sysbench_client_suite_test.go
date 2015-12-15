package sysbench_client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSysbenchClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SysbenchClient Suite")
}
