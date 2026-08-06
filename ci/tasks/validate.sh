#!/usr/bin/env bash
set -eu -o pipefail

if [[ -n "${DEBUG:-}" ]]; then
  set -x
fi

REPO_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"

"${REPO_ROOT}/scripts/lint"
echo ""
"${REPO_ROOT}/scripts/dry-run"
