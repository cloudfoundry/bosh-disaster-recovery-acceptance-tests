package testcases

import (
	_ "github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
)

type ToyTestcase struct{}

func (t ToyTestcase) Name() string {
	return "toy-testcase"
}

func (t ToyTestcase) BeforeBackup(config runner.Config) {}

func (t ToyTestcase) AfterBackup(config runner.Config) {}

func (t ToyTestcase) AfterRestore(config runner.Config) {}

func (t ToyTestcase) Cleanup(config runner.Config) {}