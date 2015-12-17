package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	It("compiles the binary", func() {
		binaryPath, err := gexec.Build("github.com/cloudfoundry-incubator/cf-mysql-benchmark-app", "-race")
		Expect(err).ToNot(HaveOccurred())
		Expect(binaryPath).To(BeAnExistingFile())
	})

	AfterEach(func() {
		gexec.CleanupBuildArtifacts()
	})
})
