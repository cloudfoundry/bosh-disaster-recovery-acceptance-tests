#!/usr/bin/env bash
set -euo pipefail

gcp_version_account_key="$( mktemp )"
echo "$GCP_SERVICE_ACCOUNT_KEY" > "$gcp_version_account_key"

function terraform_output() {
  local var="$1"
  terraform output -state=terraform-state/terraform.tfstate "$var"
}

internal_gw="$( terraform_output internal-gw )"
internal_ip="$( terraform_output jumpbox-internal-ip )"
external_ip="$( terraform_output jumpbox-ip )"
zone="$( terraform_output zone1 )"
network="$( terraform_output director-network-name )"
subnetwork="$( terraform_output director-subnetwork-name )"
tags="[$( terraform_output bosh-open-tag ), $( terraform_output jumpbox-tag )]"
internal_cidr="$( terraform_output director-subnetwork-cidr-range )"
project_id="$( terraform_output projectid )"


(
  cd jumpbox-deployment

  bosh "$BOSH_OPERATION" jumpbox.yml \
  --state "../bosh-state/${ENVIRONMENT_NAME}/state.json" \
  --vars-store "../bosh-state/${ENVIRONMENT_NAME}/creds.yml" \
  --var-file "gcp_credentials_json=${gcp_version_account_key}" \
  --ops-file gcp/cpi.yml \
  --var "internal_cidr=${internal_cidr}" \
  --var "internal_gw=${internal_gw}" \
  --var "internal_ip=${internal_ip}" \
  --var "external_ip=${external_ip}" \
  --var "tags=${tags}" \
  --var "project_id=${project_id}" \
  --var "zone=${zone}" \
  --var "network=${network}" \
  --var "subnetwork=${subnetwork}"
)

if [[ ${BOSH_OPERATION} == "delete-env" ]]; then
    rm "bosh-state/${ENVIRONMENT_NAME}/creds.yml"
fi
