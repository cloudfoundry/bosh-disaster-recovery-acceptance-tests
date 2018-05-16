package acceptance_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"testing"

	"encoding/json"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/acceptance"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	"io/ioutil"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func loadConfig() (runner.Config, runner.TestCaseFilter) {
	artifactDirPath, err := ioutil.TempDir("", "b-drats")
	if err != nil {
		panic(err)
	}

	rawConfig, err := ioutil.ReadFile(mustHaveEnv("INTEGRATION_CONFIG_PATH"))
	if err != nil {
		panic(err)
	}

	var integrationConfig acceptance.IntegrationConfig
	err = json.Unmarshal(rawConfig, &integrationConfig)
	if err != nil {
		panic(err)
	}

	config, err := runner.NewConfig(integrationConfig, mustHaveEnv("BBR_BINARY_PATH"), artifactDirPath)
	if err != nil {
		panic(err)
	}

	filter, err := runner.NewIntegrationConfigTestCaseFilter(rawConfig)
	if err != nil {
		panic(fmt.Sprint("Could not unmarshal Filter"))
	}

	return config, filter
}

func mustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	if val == "" {
		panic(fmt.Sprintf("Env var %s not set", keyname))
	}
	return val
}
