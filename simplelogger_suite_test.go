package simplelogger_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSimplelogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simplelogger Suite")
}
