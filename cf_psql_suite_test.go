package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCfPsql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cf-psql Suite")
}
