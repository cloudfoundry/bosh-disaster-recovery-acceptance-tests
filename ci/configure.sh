#!/usr/bin/env bash

set -eu

fly -t bosh-ecosystem set-pipeline -p b-drats -c ci/pipelines/b-drats/pipeline.yml
