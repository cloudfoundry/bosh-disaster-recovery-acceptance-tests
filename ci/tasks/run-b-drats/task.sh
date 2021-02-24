#!/usr/bin/env bash

set -eu

jumpbox_ip="$( terraform output -state terraform-state/terraform.tfstate jumpbox-ip | jq -r .)"
jumpbox_private_key="$( mktemp )"
bosh interpolate --path /jumpbox_ssh/private_key "bosh-vars-store/${JUMPBOX_VARS_STORE_PATH}" > "$jumpbox_private_key"
chmod 600 "$jumpbox_private_key"

export BBR_BINARY_PATH="$( pwd )/$( ls bbr-binary-release/bbr-1*-linux-amd64 )"
chmod +x "$BBR_BINARY_PATH"

export GOPATH="$( pwd )"
export PATH="${PATH}:${GOPATH}/bin"
export INTEGRATION_CONFIG_PATH="$( pwd )/b-drats-integration-config/${INTEGRATION_CONFIG_PATH}"

eval "$( ssh-agent )"
ssh-add "$jumpbox_private_key"
sshuttle -r "jumpbox@${jumpbox_ip}" "10.0.0.0/16" -D --pidfile=sshuttle.pid -e "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o ServerAliveInterval=${SSH_ALIVE_INTERVAL}"
sshuttle_pid="$( cat sshuttle.pid )"

sleep 5

if ! stat sshuttle.pid > /dev/null 2>&1; then
  echo "Failed to start sshuttle daemon"
  exit 1
fi

trap "kill ${sshuttle_pid}" EXIT


./src/github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/scripts/_run_acceptance_tests.sh

