package testcases

import (
	"fmt"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo/v2" //nolint:staticcheck
	. "github.com/onsi/gomega"    //nolint:staticcheck
)

type DeploymentTestcase struct{}

func (t DeploymentTestcase) Name() string {
	return "deployment_testcase"
}

func (t DeploymentTestcase) BeforeBackup(config runner.Config) {
	By("uploading stemcell", func() {
		runner.RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh upload stemcell",
			config,
			"upload-stemcell",
			config.StemcellSrc,
		)
	})

	By("deploying sdk deployment ", func() {
		runner.RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh deploy sdk",
			config,
			"-n",
			"-d",
			"small-deployment",
			"deploy",
			fmt.Sprintf("--var=vm_type='%s'", config.BOSH.CloudConfig.DefaultVMType),
			fmt.Sprintf("--var=network_name='%s'", config.BOSH.CloudConfig.DefaultNetwork),
			fmt.Sprintf("--var=az_name='%s'", config.BOSH.CloudConfig.DefaultAZ),
			fmt.Sprintf("--var=stemcell_os='%s'", config.StemcellOs),
			fixtures.Path("small-deployment.yml"),
		)
	})
}

func (t DeploymentTestcase) AfterBackup(config runner.Config) {
	By("deleting sdk deployment ", func() {
		runner.RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"-d",
			"small-deployment",
			"delete-deployment",
		)
	})

	By("cleaning up", func() {
		runner.RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"clean-up",
			"--all",
		)
	})
}

func (t DeploymentTestcase) AfterRestore(config runner.Config) {
	By("re-uploading stemcell", func() {
		runner.RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh upload stemcell",
			config,
			"upload-stemcell",
			config.StemcellSrc,
			"--fix",
		)
	})

	By("doing cck to bring back instances", func() {
		runner.RunBoshCommandSuccessfullyWithFailureMessage("bosh cck sdk deployment",
			config,
			"-n",
			"-d",
			"small-deployment",
			"cck",
			"--auto",
		)
	})

	By("validate deployment instances are back", func() {
		instanceState := getInstances("small-deployment", config)

		for retries := 1; !isRunning(instanceState) && retries <= 3; retries++ {
			instanceState = getInstances("small-deployment", config)
			fmt.Printf("Get instances retry attempt %v\n", retries)
			time.Sleep(time.Duration(retries*30) * time.Second)
		}

		Expect(instanceState).To(MatchRegexp("small-job/[a-z0-9-]+[ \t]+running"))
	})
}

func isRunning(state string) bool {
	return strings.Contains(state, "running")
}

func getInstances(deployment string, config runner.Config) string {
	session := runner.RunBoshCommandSuccessfullyWithFailureMessage("bosh get sdk instances",
		config,
		"-n",
		"-d",
		deployment,
		"instances",
	)
	return string(session.Out.Contents())
}

func (t DeploymentTestcase) Cleanup(config runner.Config) {
	By("deleting sdk deployment ", func() {
		runner.RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"-d",
			"small-deployment",
			"delete-deployment",
			"--force",
		)
	})
}
