package acceptance_test

import (
	"log"
	"os"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("backing up bosh", func() {
	config, filter := loadConfig()

	SetDefaultEventuallyTimeout(config.Timeout)

	testCases := []runner.TestCase{
		testcases.DeploymentTestcase{},
		testcases.TruncateDBBlobstoreTestcase{},
	}

	filteredTestCases, err := filter.Filter(testCases)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	runner.RunBoshDisasterRecoveryAcceptanceTestsSerially(config, filteredTestCases)

	AfterSuite(func() {
		By("Cleanup bosh ssh private key", func() {
			err := os.Remove(config.BOSH.SSHPrivateKeyPath)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
