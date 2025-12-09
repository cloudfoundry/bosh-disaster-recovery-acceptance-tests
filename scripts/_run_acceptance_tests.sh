#!/usr/bin/env bash
set -eu -o pipefail

SCRIPTS_DIR="$(dirname "$0")"
pushd "${SCRIPTS_DIR}/.."
  go run github.com/onsi/ginkgo/v2/ginkgo -v --trace --timeout "${GINKGO_TIMEOUT:-2h0m0s}" "acceptance"
popd
