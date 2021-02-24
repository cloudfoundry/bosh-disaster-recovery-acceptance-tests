#!/usr/bin/env bash
set -euo pipefail

jumpbox_private_key="$( mktemp )"
bosh interpolate --path /jumpbox_ssh/private_key "bosh-vars-store/${JUMPBOX_VARS_STORE_PATH}" > "$jumpbox_private_key"

jumpbox_ip="$( terraform output -state=terraform-state/terraform.tfstate jumpbox-ip | jq -r .)"
bosh_host="$( terraform output -state=terraform-state/terraform.tfstate director-internal-ip | jq -r .)"

bosh_ca_cert_path="$( mktemp )"
bosh int --path=/director_ssl/ca "bosh-vars-store/${BOSH_VARS_STORE_PATH}" > "$bosh_ca_cert_path"
bosh_client_secret="$( bosh int --path=/admin_password "bosh-vars-store/${BOSH_VARS_STORE_PATH}" )"

export BOSH_ALL_PROXY="ssh+socks5://jumpbox@${jumpbox_ip}:22?private-key=${jumpbox_private_key}"

bosh --environment "$bosh_host" \
  --client "$BOSH_CLIENT" \
  --client-secret "$bosh_client_secret" \
  --ca-cert "$bosh_ca_cert_path" \
  update-resurrection "$RESURRECTION" -n
