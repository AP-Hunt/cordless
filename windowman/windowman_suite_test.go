package windowman_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWindowman(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Windowman Suite")
}
