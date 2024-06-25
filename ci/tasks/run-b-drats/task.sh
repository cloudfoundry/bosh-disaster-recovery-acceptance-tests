#!/usr/bin/env bash

set -eu

: "${GINKGO_TIMEOUT:?}"

export GINKGO_TIMEOUT

export INTEGRATION_CONFIG_PATH="$PWD/b-drats-integration-config/${INTEGRATION_CONFIG_PATH}"

BBR_BINARY_PATH="$(ls bbr-binary-release/bbr-[0-9]*-linux-amd64)"
chmod +x "$BBR_BINARY_PATH"
export BBR_BINARY_PATH

if cat $INTEGRATION_CONFIG_PATH | jq .jumpbox_host; then
  USER=$( cat $INTEGRATION_CONFIG_PATH | jq .jumpbox_user -r )
  HOST=$( cat $INTEGRATION_CONFIG_PATH | jq .jumpbox_host -r )
  BOSH_HOST=$( cat $INTEGRATION_CONFIG_PATH | jq .bosh_host -r )
  PK_PATH=$(mktemp)
  echo -e "$( cat $INTEGRATION_CONFIG_PATH | jq .jumpbox_privkey -r)" > ${PK_PATH}

  export CREDHUB_PROXY="ssh+socks5://${USER}@${HOST}:22?private-key=${PK_PATH}"

  cat << EOF > ~/.ssh/config
Host jumphost
  HostName ${HOST}
  User ${USER}
  IdentityFile ${PK_PATH}

  StrictHostKeyChecking no
### Second jumphost. Only reachable via jumphost1.example.org
Host ${BOSH_HOST}
  HostName ${BOSH_HOST}
  ProxyJump jumphost
EOF
fi

./bosh-disaster-recovery-acceptance-tests/scripts/_run_acceptance_tests.sh
