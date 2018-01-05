#!/bin/bash

[ "$VERBOSE" ] && { set -x; export BOSH_LOG_LEVEL=debug; }
set -eu

deployment="system-test-$RANDOM"
cleanup() {
  status=$?
  ./cup --non-interactive destroy $deployment
  exit $status
}
trap cleanup EXIT

cp "$BINARY_PATH" ./cup
chmod +x ./cup

echo "DEPLOY WITH AUTOGENERATED CERT, NO DOMAIN, CUSTOM REGION, DEFAULT WORKERS"

./cup deploy $deployment

sleep 60

config=$(./cup info --json $deployment)
domain=$(echo "$config" | jq -r '.config.domain')
username=$(echo "$config" | jq -r '.config.concourse_username')
password=$(echo "$config" | jq -r '.config.concourse_password')
credhub_user=$(echo "$config" | jq -r '.secrets.credhub_username')
credhub_password=$(echo "$config" | jq -r '.secrets.credhub_password')
credhub_ca_cert=$(echo "$config" | jq -r '.secrets.credhub_ca_cert')
credhub_server="https://$(echo "$config" | jq -r '.terraform.atc_public_ip.value'):8844"
echo "$config" | jq -r '.config.concourse_ca_cert' > generated-ca-cert.pem

credhub login -u "$credhub_user" -p "$credhub_password" --ca-cert "$credhub_ca_cert" -s "$credhub_server"
credhub set -n /concourse/main/password -t password -v c1oudc0w

fly --target system-test login \
  --ca-cert generated-ca-cert.pem \
  --concourse-url "https://$domain" \
  --username "$username" \
  --password "$password"

fly --target system-test sync

fly --target system-test set-pipeline \
  --non-interactive \
  --pipeline credhub \
  --config "$(dirname "$0")/credhub.yml"

fly --target system-test unpause-pipeline \
    --pipeline credhub

fly --target system-test trigger-job \
  --job credhub/credhub \
  --watch