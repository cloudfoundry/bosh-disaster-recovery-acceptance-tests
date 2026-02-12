package testcases

import (
	"fmt"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo/v2" //nolint:staticcheck
	. "github.com/onsi/gomega"    //nolint:staticcheck
)

type TruncateDBBlobstoreTestcase struct{}

func (t TruncateDBBlobstoreTestcase) Name() string {
	return "truncate_db_blobstore_testcase"
}

func (t TruncateDBBlobstoreTestcase) BeforeBackup(config runner.Config) {
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
			fmt.Sprintf("--var=vm_type=%s", config.BOSH.CloudConfig.DefaultVMType),
			fmt.Sprintf("--var=network_name=%s", config.BOSH.CloudConfig.DefaultNetwork),
			fmt.Sprintf("--var=az_name=%s", config.BOSH.CloudConfig.DefaultAZ),
			fmt.Sprintf("--var=stemcell_os='%s'", config.StemcellOs),
			fixtures.Path("small-deployment.yml"),
		)
	})
}

func (t TruncateDBBlobstoreTestcase) AfterBackup(config runner.Config) {

	monitStop(config, "director")
	monitStop(config, "uaa")
	monitStop(config, "credhub")
	monitStop(config, "blobstore_nginx")
	monitStop(config, "postgres")

	runner.RunCommandInDirectorVMSuccessfullyWithFailureMessage(
		"truncate db/blobstore",
		config,
		"sudo rm -rf /var/vcap/store/{blobstore,director,postgres*}",
	)

	runner.RunCommandInDirectorVMSuccessfullyWithFailureMessage(
		"pre-start all jobs",
		config,
		"sudo bash -c",
		"'for pre in $(ls /var/vcap/jobs/**/bin/pre-start| grep -v monit); do $pre; done'",
	)

	fmt.Println("Waiting for  all jobs to finish pre-start")
	runner.RunCommandInDirectorVMSuccessfullyWithFailureMessage(
		"pre-start all jobs",
		config,
		"sudo bash -c",
		`'while pgrep pre-start; do sleep 1; echo "waiting for pre-start to finish"; done'`,
	)

	monitStart(config, "postgres")
	monitStart(config, "blobstore_nginx")
	monitStart(config, "uaa")
	monitStart(config, "credhub")
	monitStart(config, "director")

	Eventually(func() int {
		session := runner.RunBoshCommand(
			"bosh env",
			config,
			"env",
		)
		return session.Wait().ExitCode()
	}).Should(Equal(0))
}

func (t TruncateDBBlobstoreTestcase) AfterRestore(config runner.Config) {
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
		Eventually(func() string {
			session := runner.RunBoshCommandSuccessfullyWithFailureMessage("bosh get sdk instances",
				config,
				"-n",
				"-d",
				"small-deployment",
				"instances",
			)
			return string(session.Out.Contents())
		}, fixtures.EventuallyTimeout, fixtures.EventuallyRetryInterval).Should(MatchRegexp("small-job/[a-z0-9-]+[ \t]+running"))
	})
}

func (t TruncateDBBlobstoreTestcase) Cleanup(config runner.Config) {
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
