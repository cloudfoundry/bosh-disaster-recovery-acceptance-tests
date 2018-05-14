package runner

import (
	"fmt"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunBoshDisasterRecoveryAcceptanceTests(config Config, testCases []TestCase) {
	It("backs up and restores bosh", func() {
		By("running the before backup step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the before backup step for " + testCase.Name())
				testCase.BeforeBackup(config)
			}
		})

		By("backing up", func() {
			RunCommandSuccessfullyWithFailureMessage(
				"bbr director backup",
				os.Stdout,
				fmt.Sprintf(
					"%s director --host %s --username %s --private-key-path %s backup --artifact-path %s",
					config.BBRBinaryPath,
					config.BOSH.Host,
					config.BOSH.SSHUsername,
					config.BOSH.SSHPrivateKeyPath,
					config.ArtifactPath,
				),
			)
		})

		By("running the after backup step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the after backup step for " + testCase.Name())
				testCase.AfterBackup(config)
			}
		})

		By("restoring", func() {
			RunCommandSuccessfullyWithFailureMessage(
				"bbr director restore",
				os.Stdout,
				fmt.Sprintf(
					"%s director --host %s --username %s --private-key-path %s "+
						"restore --artifact-path %s/$(ls %s | grep %s | head -n 1)",
					config.BBRBinaryPath,
					config.BOSH.Host,
					config.BOSH.SSHUsername,
					config.BOSH.SSHPrivateKeyPath,
					config.ArtifactPath,
					config.ArtifactPath,
					config.BOSH.Host,
				),
			)
		})

		By("waiting for bosh director api to be available", func() {
			Eventually(func() int {
				return RunBoshCommand("bosh releases", GinkgoWriter, config, "releases").ExitCode()
			}, "60s", "1s").Should(BeZero())
		})

		By("running the after restore step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the after restore step for " + testCase.Name())
				testCase.AfterRestore(config)
			}
		})
	})

	AfterEach(func() {
		By("bbr director backup-cleanup", func() {
			RunCommandSuccessfullyWithFailureMessage(
				"bbr director backup-cleanup",
				os.Stdout,
				fmt.Sprintf(
					"%s director --host %s --username %s --private-key-path %s backup-cleanup",
					config.BBRBinaryPath,
					config.BOSH.Host,
					config.BOSH.SSHUsername,
					config.BOSH.SSHPrivateKeyPath,
				),
			)
		})

		By("Running cleanup for each testcase", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the cleanup step for " + testCase.Name())
				testCase.Cleanup(config)
			}
		})

		By("Cleanup bosh ssh private key", func() {
			err := os.Remove(config.BOSH.SSHPrivateKeyPath)
			Expect(err).NotTo(HaveOccurred())
		})

		By("Cleanup bbr director backup artifact", func() {
			err := os.RemoveAll(config.ArtifactPath)
			Expect(err).NotTo(HaveOccurred())
		})
	})
}
