package urls_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestURLs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "URLs Suite")
}
