package acceptance_test

import (
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("backing up bosh", func() {
	config := loadConfig()

	testCases := []runner.TestCase{
		testcases.DeploymentTestcase{},
	}

	SetDefaultEventuallyTimeout(config.Timeout)
	runner.RunBoshDisasterRecoveryAcceptanceTests(config, testCases)
})
