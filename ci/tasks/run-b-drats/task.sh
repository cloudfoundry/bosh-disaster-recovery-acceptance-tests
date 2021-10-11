#!/usr/bin/env bash

set -eu

: "${JUMPBOX_IP:="$( terraform output -state terraform-state/terraform.tfstate jumpbox-ip | jq -r .)"}"
: "${JUMPBOX_PRIVATE_KEY:="$( bosh interpolate --path /jumpbox_ssh/private_key "bosh-vars-store/${JUMPBOX_VARS_STORE_PATH}" )"}"
: "${JUMPBOX_USER:=jumpbox}"

jumpbox_private_key="$( mktemp )"
echo -e "$JUMPBOX_PRIVATE_KEY" | sed -e 's/^"//' -e 's/"$//' > "$jumpbox_private_key"
chmod 600 "$jumpbox_private_key"

eval "$( ssh-agent )"
ssh-add "$jumpbox_private_key"
sshuttle -r "${JUMPBOX_USER}@${JUMPBOX_IP}" "10.0.0.0/16" -D --pidfile=sshuttle.pid -e "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o ServerAliveInterval=${SSH_ALIVE_INTERVAL}"
sshuttle_pid="$( cat sshuttle.pid )"

sleep 5

if ! stat sshuttle.pid > /dev/null 2>&1; then
  echo "Failed to start sshuttle daemon"
  exit 1
fi

trap 'kill ${sshuttle_pid}' EXIT

BBR_BINARY_PATH="${PWD}/$( ls bbr-binary-release/bbr-1*-linux-amd64 )"
chmod +x "$BBR_BINARY_PATH"
export BBR_BINARY_PATH

GOPATH="$( pwd )"
PATH="${PATH}:${GOPATH}/bin"
INTEGRATION_CONFIG_PATH="$( pwd )/b-drats-integration-config/${INTEGRATION_CONFIG_PATH}"

export GOPATH PATH INTEGRATION_CONFIG_PATH

./bosh-disaster-recovery-acceptance-tests/scripts/_run_acceptance_tests.sh

