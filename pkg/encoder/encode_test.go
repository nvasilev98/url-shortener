package encoder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"url-shortener/pkg/encoder"
)

var _ = Describe("Encoder", func() {
	When("encoding a number to base62", func() {
		It("should return base62 interpretation of the number", func() {
			encodedNumber := encoder.New().EncodeToBase62(3256)
			Expect(encodedNumber).To(Equal("qW"))
		})
	})
})
