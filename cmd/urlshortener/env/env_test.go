package env_test

import (
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"url-shortener/cmd/urlshortener/env"
)

var _ = Describe("Env", func() {
	const (
		hostEnv     = "HOST"
		portEnv     = "PORT"
		host        = "127.0.0.1"
		port        = 8080
		invalidPort = "invalid"
	)

	When("env variables are set incorrectly", func() {
		BeforeEach(func() {
			Expect(os.Setenv(hostEnv, host)).To(Succeed())
			Expect(os.Setenv(portEnv, invalidPort)).To(Succeed())
		})

		AfterEach(func() {
			Expect(os.Unsetenv(hostEnv)).To(Succeed())
			Expect(os.Unsetenv(portEnv)).To(Succeed())
		})

		It("should return an error", func() {
			_, err := env.LoadAppConfig()
			Expect(err).To(HaveOccurred())
		})
	})

	When("env variables are set correctly", func() {
		BeforeEach(func() {
			Expect(os.Setenv(hostEnv, host)).To(Succeed())
			Expect(os.Setenv(portEnv, strconv.Itoa(port))).To(Succeed())
		})

		AfterEach(func() {
			Expect(os.Unsetenv(hostEnv)).To(Succeed())
			Expect(os.Unsetenv(portEnv)).To(Succeed())
		})

		It("should load app config", func() {
			config, err := env.LoadAppConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config.Host).To(Equal(host))
			Expect(config.Port).To(Equal(port))
		})
	})
})
