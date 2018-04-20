#!/usr/bin/env bash

set -eu

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

pushd src/github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests
  scripts/_run_acceptance_tests.sh
popd
