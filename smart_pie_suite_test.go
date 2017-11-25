package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

const (
	BROKER_URI = "tcp://localhost:55001"
)

func TestSmartPie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SmartPie Suite")
}

var _ = BeforeSuite(func() {
})

var _ = AfterSuite(func() {
})
