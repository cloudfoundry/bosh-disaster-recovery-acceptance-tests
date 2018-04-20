#!/usr/bin/env bash

set -eu

go get github.com/onsi/ginkgo/ginkgo
dep ensure
ginkgo -v --trace acceptance
