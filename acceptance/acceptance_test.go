package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/testcases"
)

var _ = Describe("backing up bosh", func() {
	testCases := []runner.TestCase{
		testcases.ToyTestcase{},
	}

	runner.RunBoshDisasterRecoveryAcceptanceTests(testCases)
})
