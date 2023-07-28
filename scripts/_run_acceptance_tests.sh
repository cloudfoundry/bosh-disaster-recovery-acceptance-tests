#!/usr/bin/env bash

set -eu

SCRIPTS_DIR="$(dirname "$0")"
pushd "$SCRIPTS_DIR/.."
  go mod tidy
  go run github.com/onsi/ginkgo/v2/ginkgo -v --trace --timeout "$GINKGO_TIMEOUT" "acceptance"
popd
