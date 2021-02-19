#!/usr/bin/env bash

set -eu

SCRIPTS_DIR="$(dirname $0)"
pushd "$SCRIPTS_DIR/.."
  go mod tidy
  ginkgo -v --trace "acceptance"
popd
