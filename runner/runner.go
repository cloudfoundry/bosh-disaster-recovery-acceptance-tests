package runner

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

func RunBoshDisasterRecoveryAcceptanceTests(config Config, testCases []TestCase) {
	It("backs up and restores bosh", func() {
		By("running the before backup step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the before backup step for " + testCase.Name())
				testCase.BeforeBackup()
			}
		})

		By("backing up", func() {
			RunCommandSuccessfullyWithFailureMessage(
				"bbr director backup",
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
				testCase.AfterBackup()
			}
		})

		By("restoring", func() {
			RunCommandSuccessfullyWithFailureMessage(
				"bbr director restore",
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

		By("running the after restore step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the after restore step for " + testCase.Name())
				testCase.AfterRestore()
			}
		})
	})

	AfterEach(func() {
		By("bbr director backup-cleanup", func() {
			RunCommandSuccessfullyWithFailureMessage(
				"bbr director backup-cleanup",
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
				testCase.Cleanup()
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
