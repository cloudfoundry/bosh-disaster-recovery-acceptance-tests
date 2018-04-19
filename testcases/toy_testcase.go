package testcases

import (
	_ "github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
)

type ToyTestcase struct{}

func (t ToyTestcase) Name() string {
	return "toy-testcase"
}

func (t ToyTestcase) BeforeBackup() {}

func (t ToyTestcase) AfterBackup() {}

func (t ToyTestcase) AfterRestore() {}

func (t ToyTestcase) Cleanup() {}