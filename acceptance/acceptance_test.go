package acceptance_test

import (
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"os"
)

var _ = Describe("backing up bosh", Ordered, func() {
	config, filter := loadConfig()

	SetDefaultEventuallyTimeout(config.Timeout)

	testCases := []runner.TestCase{
		testcases.DeploymentTestcase{},
		testcases.TruncateDBBlobstoreTestcase{},
		testcases.CredhubTestcase{},
	}

	filteredTestCases, err := filter.Filter(testCases)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	runner.RunBoshDisasterRecoveryAcceptanceTestsSerially(config, filteredTestCases)

	AfterAll(func() {
		By("Cleanup bosh ssh private key", func() {
			err := os.Remove(config.BOSH.SSHPrivateKeyPath)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
