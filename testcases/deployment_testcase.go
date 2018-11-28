package testcases

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	. "github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DeploymentTestcase struct{}

func (t DeploymentTestcase) Name() string {
	return "deployment_testcase"
}

func (t DeploymentTestcase) BeforeBackup(config Config) {
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

func getStemcell(config Config) string {
	sess := RunBoshCommandSuccessfullyWithFailureMessage(
		"bosh stemcells",
		config,
		"stemcells",
	)
	stemcell := string(sess.Out.Contents())
	if strings.Contains(stemcell, "trusty") {
		stemcell = "trusty"
	} else {
		stemcell = "xenial"
	}
	return stemcell
}

func (t DeploymentTestcase) AfterBackup(config Config) {
	By("deleting sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
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
		Expect(session.Out.Contents()).To(MatchRegexp("database-backuper/[a-z0-9-]+[ \t]+running"))
	})

}

func (t DeploymentTestcase) Cleanup(config Config) {
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
