#!/usr/bin/env bash

set -eu

SCRIPTS_DIR="$(dirname $0)"

ginkgo -v --trace "$SCRIPTS_DIR/../acceptance"
