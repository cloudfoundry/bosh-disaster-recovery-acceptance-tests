#!/usr/bin/env bash
# shellcheck disable=SC2034,SC2153

set -euo pipefail

bosh_host="$(terraform output -state terraform-state/terraform.tfstate director-internal-ip | jq -r .)"
bosh_ssh_username="$BOSH_SSH_USERNAME"
bosh_ssh_private_key="$( bosh int --path=/jumpbox_ssh/private_key "director-creds/creds.yml" )"
timeout_in_minutes="$TIMEOUT_IN_MINUTES"
bosh_client="$BOSH_CLIENT"
bosh_client_secret="$( bosh int --path=/admin_password "director-creds/creds.yml" )"
bosh_ca_cert="$( bosh int --path=/director_ssl/ca "director-creds/creds.yml" )"
include_deployment_testcase="$INCLUDE_DEPLOYMENT_TESTCASE"
include_truncate_db_blobstore_testcase="$INCLUDE_TRUNCATE_DB_BLOBSTORE_TESTCASE"
include_credhub_testcase="$INCLUDE_CREDHUB_TESTCASE"
credhub_client="$CREDHUB_CLIENT"
credhub_client_secret="$( bosh int --path=/credhub_admin_client_secret "director-creds/creds.yml" )"
credhub_server="$CREDHUB_SERVER"
credhub_ca_cert="$( bosh interpolate "director-creds/creds.yml" --path=/credhub_tls/ca )
$( bosh interpolate "director-creds/creds.yml" --path=/uaa_ssl/ca )"
stemcell_src="$( cat stemcell/url )"
jumpbox_host="$(terraform output -state terraform-state/terraform.tfstate jumpbox-ip | jq -r .)"
jumpbox_user="jumpbox"
jumpbox_pubkey="$(bosh interpolate --path /jumpbox_ssh/public_key "jumpbox-creds/creds.yml")"
jumpbox_privkey="$(bosh interpolate --path /jumpbox_ssh/private_key "jumpbox-creds/creds.yml")"

integration_config="{}"

string_vars="bosh_host bosh_ssh_username bosh_ssh_private_key bosh_client bosh_client_secret bosh_ca_cert stemcell_src credhub_client_secret credhub_client credhub_ca_cert credhub_server jumpbox_host jumpbox_user jumpbox_pubkey jumpbox_privkey"
for var in $string_vars
do
  integration_config=$(echo ${integration_config} | jq ".${var}=\"${!var}\"")
done

other_vars="include_deployment_testcase include_truncate_db_blobstore_testcase include_credhub_testcase timeout_in_minutes"
for var in $other_vars
do
  integration_config=$(echo "${integration_config}" | jq ".${var}=${!var}")
done

echo "$integration_config" > b-drats-integration-config/integration_config.json
