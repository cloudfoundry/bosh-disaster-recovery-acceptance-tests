package testcases

import (
	"fmt"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	. "github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/gomega"
)

func monitWait(config Config, jobName, status string) {
	Eventually(func() string {
		statusLines := strings.Split(string(RunCommandInDirectorVM(
			"monit summary",
			config,
			"sudo /var/vcap/bosh/bin/monit summary",
		).Wait().Out.Contents()), "\n")

		for _, statusLine := range statusLines {
			if strings.Contains(statusLine, fmt.Sprintf("Process '%s'", jobName)) {
				return statusLine
			}
		}
		return ""
	}, fixtures.EventuallyTimeout, fixtures.EventuallyRetryInterval).Should(HaveSuffix(status))
}

func monitStop(config Config, jobName string) {
	RunCommandInDirectorVMSuccessfullyWithFailureMessage(
		fmt.Sprintf("monit stop %s", jobName),
		config,
		fmt.Sprintf("sudo /var/vcap/bosh/bin/monit stop %s", jobName),
	)

	monitWait(config, jobName, "not monitored")
}

func monitStart(config Config, jobName string) {
	RunCommandInDirectorVMSuccessfullyWithFailureMessage(
		fmt.Sprintf("monit start %s", jobName),
		config,
		fmt.Sprintf("sudo /var/vcap/bosh/bin/monit start %s", jobName),
	)

	monitWait(config, jobName, "running")
}

func attemptWithBackoff(function func() error, attempts int) error {
	var err error
	err = function()
	if err == nil {
		return nil
	}

	for i := 1; i < attempts; i++ {
		backoff := time.Duration(i)
		time.Sleep(5 * backoff * time.Second)
		err = function()
		if err == nil {
			return nil
		}
	}
	return err
}
