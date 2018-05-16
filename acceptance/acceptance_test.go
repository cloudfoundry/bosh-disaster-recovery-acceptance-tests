package acceptance_test

import (
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

var _ = Describe("backing up bosh", func() {
	config, filter := loadConfig()

	SetDefaultEventuallyTimeout(config.Timeout)

	testCases := []runner.TestCase{
		testcases.DeploymentTestcase{},
	}

	filteredTestCases, err := filter.Filter(testCases)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	runner.RunBoshDisasterRecoveryAcceptanceTests(config, filteredTestCases)
})
