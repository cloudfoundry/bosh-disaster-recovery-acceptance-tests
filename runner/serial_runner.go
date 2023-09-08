package runner

import (
	"fmt"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path"
)

func RunBoshDisasterRecoveryAcceptanceTestsSerially(config Config, testCases []TestCase) {
	fmt.Println("Running testcases: ")
	for _, t := range testCases {
		fmt.Println(t.Name())
	}

	for _, test := range testCases {
		testCase := test

		Context(fmt.Sprintf("test case %s", testCase.Name()), func() {
			var (
				artifactPath       commonArtifactPath
				boshPrivateKeyPath string
			)

			BeforeEach(func() {
				artifactPath = newArtifactPath(config)

				// When running with a jumpbox, the private key is already copied across
				boshPrivateKeyPath = config.BOSH.SSHPrivateKeyPath
				if config.Jumpbox != nil {
					boshPrivateKeyPath = path.Join("/tmp", path.Base(config.BOSH.SSHPrivateKeyPath))
				}
			})

			AfterEach(func() {
				By("bbr director backup-cleanup", func() {
					RunBBRCommandSuccessfullyWithFailureMessage(
						"bbr director backup-cleanup",
						config,
						fmt.Sprintf(
							"director --host %s --username %s --private-key-path %s backup-cleanup",
							config.BOSH.Host,
							config.BOSH.SSHUsername,
							boshPrivateKeyPath,
						),
					)
				})

				By("Running cleanup for each testcase", func() {
					fmt.Println("Running the cleanup step for " + testCase.Name())
					testCase.Cleanup(config)
				})

				By("Cleanup bbr director backup artifact", func() {
					artifactPath.cleanup()
				})
			})

			It("backs up and restores bosh", func() {
				By("running the before backup step", func() {
					fmt.Println("Running the before backup step for " + testCase.Name())
					testCase.BeforeBackup(config)
				})

				By("backing up", func() {
					RunBBRCommandSuccessfullyWithFailureMessage(
						"bbr director backup",
						config,
						fmt.Sprintf(
							"director --host %s --username %s --private-key-path %s backup --artifact-path %s",
							config.BOSH.Host,
							config.BOSH.SSHUsername,
							boshPrivateKeyPath,
							artifactPath.path(),
						),
					)
				})

				By("running the after backup step", func() {
					fmt.Println("Running the after backup step for " + testCase.Name())
					testCase.AfterBackup(config)
				})

				artifactToRestore := artifactPath.firstMatch(config.BOSH.Host) // prints output

				By("restoring", func() {
					RunBBRCommandSuccessfullyWithFailureMessage(
						"bbr director restore",
						config,
						fmt.Sprintf(
							"director --host %s --username %s --private-key-path %s "+
								"restore --artifact-path %s",
							config.BOSH.Host,
							config.BOSH.SSHUsername,
							boshPrivateKeyPath,
							artifactToRestore,
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

	BeforeAll(func() {
		config.Jumpbox.Deploy(config)
	})

	AfterAll(func() {
		config.Jumpbox.Cleanup(config)
	})
}
