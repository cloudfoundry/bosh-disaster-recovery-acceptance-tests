#!/usr/bin/env bash
set -euo pipefail

jumpbox_private_key="$( mktemp )"
bosh interpolate --path /jumpbox_ssh/private_key "jumpbox-creds/creds.yml" > "$jumpbox_private_key"

jumpbox_ip="$( terraform output -state=terraform-state/terraform.tfstate jumpbox-ip | jq -r .)"
bosh_host="$( terraform output -state=terraform-state/terraform.tfstate director-internal-ip | jq -r .)"

bosh_ca_cert_path="$( mktemp )"
bosh int --path=/director_ssl/ca "director-creds/creds.yml" > "$bosh_ca_cert_path"
bosh_client_secret="$( bosh int --path=/admin_password "director-creds/creds.yml" )"

export BOSH_ALL_PROXY="ssh+socks5://jumpbox@${jumpbox_ip}:22?private-key=${jumpbox_private_key}"

bosh --environment "$bosh_host" \
  --client "$BOSH_CLIENT" \
  --client-secret "$bosh_client_secret" \
  --ca-cert "$bosh_ca_cert_path" \
  update-resurrection "$RESURRECTION" -n
