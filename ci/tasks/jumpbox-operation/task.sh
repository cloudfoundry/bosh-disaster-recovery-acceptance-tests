#!/usr/bin/env bash
set -euo pipefail

gcp_version_account_key="$( mktemp )"
echo "$GCP_SERVICE_ACCOUNT_KEY" > "$gcp_version_account_key"

function terraform_output() {
  local var="$1"
  jq -r ".\"${var}\"" terraform-state/metadata
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

pushd jumpbox-deployment
  bosh "$BOSH_OPERATION" jumpbox.yml \
    --state "../jumpbox-state/state.json" \
    --vars-store "../jumpbox-creds/creds.yml" \
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
popd

if [[ ${BOSH_OPERATION} == "delete-env" ]]; then
    echo "" > "jumpbox-creds/creds.yml"
    echo "{}" > "jumpbox-state/state.json"
fi
