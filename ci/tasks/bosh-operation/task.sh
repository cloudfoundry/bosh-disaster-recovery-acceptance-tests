#!/usr/bin/env bash
set -euo pipefail

gcp_version_account_key="$(mktemp)"
echo "${GCP_SERVICE_ACCOUNT_KEY}" > "${gcp_version_account_key}"

jumpbox_private_key="$(mktemp)"
bosh interpolate --path /jumpbox_ssh/private_key "jumpbox-creds/creds.yml" > "$jumpbox_private_key"

function terraform_output() {
  local var="$1"
  jq -r ".\"${var}\"" terraform-state/metadata
}

jumpbox_ip="$(terraform_output jumpbox-ip)"
internal_gw="$(terraform_output internal-gw)"
internal_ip="$(terraform_output director-internal-ip)"
zone="$(terraform_output zone1)"
network="$(terraform_output director-network-name)"
subnetwork="$(terraform_output director-subnetwork-name)"
tags="[$(terraform_output director-tag)]"
internal_cidr="$(terraform_output director-subnetwork-cidr-range)"
project_id="$(terraform_output projectid)"

export BOSH_ALL_PROXY="ssh+socks5://jumpbox@${jumpbox_ip}:22?private-key=${jumpbox_private_key}"

pushd bosh-deployment
  bosh "$BOSH_OPERATION" bosh.yml \
  --state "../director-state/state.json" \
  --vars-store "../director-creds/creds.yml" \
  --var-file gcp_credentials_json="${gcp_version_account_key}" \
  --ops-file gcp/cpi.yml \
  --ops-file uaa.yml \
  --ops-file credhub.yml \
  --ops-file jumpbox-user.yml \
  --ops-file bbr.yml \
  --ops-file ../bosh-disaster-recovery-acceptance-tests/ci/infrastructure/opsfiles/gcp/bosh-director-ephemeral-ip-ops.yml \
  --var "director_name=${DIRECTOR_NAME}" \
  --var "internal_cidr=${internal_cidr}" \
  --var "internal_gw=${internal_gw}" \
  --var "internal_ip=${internal_ip}" \
  --var "project_id=${project_id}" \
  --var "zone=${zone}" \
  --var "tags=${tags}" \
  --var "network=${network}" \
  --var "subnetwork=${subnetwork}"
popd

if [[ ${BOSH_OPERATION} == "delete-env" ]]; then
  echo "" > "director-creds/creds.yml"
  echo "{}" > "director-state/state.json"
fi

