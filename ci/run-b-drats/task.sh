#!/usr/bin/env bash

set -e
set -u

pushd $BBR_RELEASE_PATH
  tar xvf *.tar
  export BBR_BINARY_PATH="$PWD/releases/bbr"
popd

export GOPATH="$PWD"
export PATH="$PATH:$GOPATH/bin"
INTEGRATION_CONFIG_PATH="$PWD/${INTEGRATION_CONFIG_PATH}"

./$B_DRATS_PATH/scripts/_run_acceptance_tests.sh

