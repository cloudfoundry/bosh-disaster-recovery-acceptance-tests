package runner

import (
	"fmt"
	"io/ioutil"

	"os"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunBoshDisasterRecoveryAcceptanceTestsSerially(config Config, testCases []TestCase) {
	fmt.Println("Running testcases: ")
	for _, t := range testCases {
		fmt.Println(t.Name())
	}

	for _, test := range testCases {
		testCase := test

		Context(fmt.Sprintf("test case %s", testCase.Name()), func() {
			var artifactPath string
			BeforeEach(func() {
				var err error
				artifactPath, err = ioutil.TempDir("", "b-drats")
				Expect(err).NotTo(HaveOccurred())
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
					fmt.Println("Running the cleanup step for " + testCase.Name())
					testCase.Cleanup(config)
				})

				By("Cleanup bbr director backup artifact", func() {
					err := os.RemoveAll(artifactPath)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			It("backs up and restores bosh", func() {
				By("running the before backup step", func() {
					fmt.Println("Running the before backup step for " + testCase.Name())
					testCase.BeforeBackup(config)
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
							artifactPath,
						),
					)
				})

				By("running the after backup step", func() {
					fmt.Println("Running the after backup step for " + testCase.Name())
					testCase.AfterBackup(config)
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
							artifactPath,
							artifactPath,
							config.BOSH.Host,
						),
					)
				})

				By("waiting for bosh director api to be available", func() {
					Eventually(func() int {
						return RunBoshCommand("bosh releases", config, "releases").ExitCode()
					}, fixtures.EventuallyTimeout, fixtures.EventuallyRetryInterval).Should(BeZero())
				})

				By("running the after restore step", func() {
					fmt.Println("Running the after restore step for " + testCase.Name())
					testCase.AfterRestore(config)
				})
			})
		})
	}
}
