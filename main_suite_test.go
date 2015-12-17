package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBenchmarkConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Benchmark App Main Suite")
}
