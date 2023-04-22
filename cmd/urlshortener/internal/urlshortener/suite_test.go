package urlshortener_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestURLShortener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "URLShortener Suite")
}
