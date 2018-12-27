package testcases

import (
	"fmt"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	. "github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CleanUpTestcase struct{}

func (t CleanUpTestcase) Name() string {
	return "clean_up_testcase"
}

func (t CleanUpTestcase) BeforeBackup(config Config) {
	By("uploading stemcell", func() {
		RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh upload stemcell",
			config,
			"upload-stemcell",
			config.StemcellPath,
		)
	})

	By("deploying sdk deployment ", func() {
		stemcell := getStemcell(config)
		RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh deploy sdk",
			config,
			"-n",
			"-d",
			"sdk-test",
			"deploy",
			fmt.Sprintf("--var=vm_type=%s", config.BOSH.CloudConfig.DefaultVMType),
			fmt.Sprintf("--var=network_name=%s", config.BOSH.CloudConfig.DefaultNetwork),
			fmt.Sprintf("--var=az_name=%s", config.BOSH.CloudConfig.DefaultAZ),
			fixtures.Path(fmt.Sprintf("sdk-manifest-%s.yml", stemcell)),
		)
	})
}

func (t CleanUpTestcase) AfterBackup(config Config) {
	By("deleting sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"-d",
			"sdk-test",
			"delete-deployment",
		)
	})

	By("cleaning up", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"clean-up",
			"--all",
		)
	})
}

func (t CleanUpTestcase) AfterRestore(config Config) {
	By("re-uploading stemcell", func() {
		RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh upload stemcell",
			config,
			"upload-stemcell",
			config.StemcellPath,
			"--fix",
		)
	})

	By("doing cck to bring back instances", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh cck sdk deployment",
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
			config,
			"-n",
			"-d",
			"sdk-test",
			"instances",
		)
		Expect(string(session.Out.Contents())).To(MatchRegexp("database-backuper/[a-z0-9-]+[ \t]+running"))
	})
}

func (t CleanUpTestcase) Cleanup(config Config) {
	By("deleting sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"-d",
			"sdk-test",
			"delete-deployment",
			"--force",
		)
	})
}
