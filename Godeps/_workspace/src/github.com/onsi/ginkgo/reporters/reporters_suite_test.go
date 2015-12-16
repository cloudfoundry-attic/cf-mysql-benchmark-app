package reporters_test

import (
	. "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry-incubator/cf-mysql-benchmark-app/Godeps/_workspace/src/github.com/onsi/gomega"
	"testing"
)

func TestReporters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reporters Suite")
}
