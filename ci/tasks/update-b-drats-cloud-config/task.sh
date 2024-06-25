#!/usr/bin/env bash
set -euo pipefail

jumpbox_private_key="$( mktemp )"
bosh interpolate --path /jumpbox_ssh/private_key "jumpbox-creds/creds.yml" > "$jumpbox_private_key"

function terraform_output() {
  local var="$1"
  jq -r ".\"${var}\"" terraform-state/metadata
}

bosh_host="$( terraform_output director-internal-ip )"
jumpbox_ip="$( terraform_output jumpbox-ip )"
jumpbox_internal_ip="$( terraform_output jumpbox-internal-ip )"
internal_gw="$( terraform_output internal-gw )"
zone="$( terraform_output zone1 )"
network="$( terraform_output director-network-name )"
subnetwork="$( terraform_output director-subnetwork-name )"
tags="[$( terraform_output internal-tag )]"
internal_cidr="$( terraform_output director-subnetwork-cidr-range )"

bosh_ca_cert_path="$( mktemp )"
bosh int --path=/director_ssl/ca "director-creds/creds.yml" > "$bosh_ca_cert_path"
bosh_client_secret="$( bosh int --path=/admin_password "director-creds/creds.yml" )"

export BOSH_ALL_PROXY="ssh+socks5://jumpbox@${jumpbox_ip}:22?private-key=${jumpbox_private_key}"

bosh --environment "$bosh_host" \
  --client "$BOSH_CLIENT" \
  --client-secret "$bosh_client_secret" \
  --ca-cert "$bosh_ca_cert_path" \
  update-cloud-config \
  --ops-file bosh-disaster-recovery-acceptance-tests/ci/infrastructure/opsfiles/gcp/cloud-config-jumpbox-reserved-ip.yml \
  --var "internal_cidr=${internal_cidr}" \
  --var "internal_gw=${internal_gw}"\
  --var "zone=${zone}" \
  --var "network=${network}" \
  --var "subnetwork=${subnetwork}" \
  --var "tags=${tags}" \
  --var "subnetwork_reserved_ip=${jumpbox_internal_ip}" \
  "cloud-config/${CLOUD_CONFIG_PATH}" -n
