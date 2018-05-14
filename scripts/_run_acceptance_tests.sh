#!/usr/bin/env bash

set -eu

SCRIPTS_DIR="$(dirname $0)"

ginkgo --trace "$SCRIPTS_DIR/../acceptance"
