package testcases

import (
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	. "github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

type DeploymentTestcase struct{}

func (t DeploymentTestcase) Name() string {
	return "deployment-testcase"
}

func (t DeploymentTestcase) BeforeBackup(config Config) {
	By("deploying sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh deploy sdk",
			os.Stdout,
			config,
			"-n",
			"-d",
			"sdk-test",
			"deploy",
			fixtures.Path("sdk-manifest.yml"),
		)
	})
}

func (t DeploymentTestcase) AfterBackup(config Config) {
	By("deleting sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			os.Stdout,
			config,
			"-n",
			"-d",
			"sdk-test",
			"delete-deployment",
		)
	})
}

func (t DeploymentTestcase) AfterRestore(config Config) {
	By("doing cck to bring back instances", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh cck sdk deployment",
			os.Stdout,
			config,
			"-n",
			"-d",
			"sdk-test",
			"cck",
			"--auto",
		)
	})

	By("validate deployment instances are back", func() {
		session := RunBoshCommandSuccessfullyWithFailureMessage("bosh get sdk instances",
			os.Stdout,
			config,
			"-n",
			"-d",
			"sdk-test",
			"instances",
		)
		Expect(session.Out.Contents()).To(MatchRegexp("database-backuper/[a-z0-9-]+[ \t]+running"))
	})

}

func (t DeploymentTestcase) Cleanup(config Config) {
	By("deleting sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			os.Stdout,
			config,
			"-n",
			"-d",
			"sdk-test",
			"delete-deployment",
			"--force",
		)
	})
}
