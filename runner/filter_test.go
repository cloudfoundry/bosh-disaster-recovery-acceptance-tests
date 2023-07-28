package runner_test

import (
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestCaseFilter", func() {
	Describe("IntegrationConfigTestCaseFilter", func() {
		var filter runner.IntegrationConfigTestCaseFilter

		BeforeEach(func() {
			filter = runner.IntegrationConfigTestCaseFilter(map[string]interface{}{
				"include_one":   true,
				"include_two":   false,
				"include_three": true,
			})
		})

		It("only runs tests that are included in the config", func() {
			filteredTestCases, err := filter.Filter(testCases("one", "two", "three"))

			Expect(err).NotTo(HaveOccurred())
			Expect(filteredTestCases).To(ConsistOf(testCases("one", "three")))
		})

		Context("when not flag is specified for a testcase", func() {
			It("defaults to false", func() {
				filteredTestCases, err := filter.Filter(testCases("one", "two", "four"))

				Expect(err).NotTo(HaveOccurred())
				Expect(filteredTestCases).To(ConsistOf(testCases("one")))
			})
		})

		Context("when no test case matches", func() {
			It("returns an error", func() {
				filteredTestCases, err := filter.Filter(testCases("six"))

				Expect(filteredTestCases).To(BeNil())
				Expect(err).To(MatchError("unable to find any test case included by the config"))
			})
		})

		Context("when an include flag has a non-boolean value", func() {
			BeforeEach(func() {
				filter = runner.IntegrationConfigTestCaseFilter(map[string]interface{}{
					"include_one": true,
					"include_two": "not a boolean",
				})
			})

			It("returns an error", func() {
				filteredTestCases, err := filter.Filter(testCases("one", "two"))

				Expect(filteredTestCases).To(BeNil())
				Expect(err).To(MatchError("'include_two' should be a boolean"))
			})
		})
	})
})

type FakeTestCase struct {
	name string
}

func (tc FakeTestCase) Name() string {
	return tc.name
}

func (tc FakeTestCase) BeforeBackup(config runner.Config) {}
func (tc FakeTestCase) AfterBackup(config runner.Config)  {}
func (tc FakeTestCase) AfterRestore(config runner.Config) {}
func (tc FakeTestCase) Cleanup(config runner.Config)      {}

func testCase(name string) runner.TestCase {
	return FakeTestCase{name: name}
}

func testCases(names ...string) []runner.TestCase {
	var tcs []runner.TestCase

	for _, name := range names {
		tcs = append(tcs, testCase(name))
	}

	return tcs
}
