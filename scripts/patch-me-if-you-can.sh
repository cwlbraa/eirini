#!/bin/bash

set -xeuo pipefail

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
eirini_basedir=$(realpath "$script_dir/..")
eirini_release_basedir=$(realpath "$script_dir/../../eirini-release")

docker_build() {
  echo "Building docker image for $1"
  pushd "$eirini_basedir"
    docker build . -f "$eirini_basedir/docker/$component/Dockerfile" \
      --build-arg GIT_SHA=big-sha \
      --tag "eirini/$component:patch-me-if-you-can"
  popd
}

docker_push() {
  echo "Pushing docker image for $1"
  pushd "$eirini_basedir"
    docker push "eirini/$component:patch-me-if-you-can"
  popd
}

update_helm_chart() {
  echo "Applying docker image of $1 to kubernetes cluster"
  pushd "$eirini_release_basedir/helm/eirini/templates"
    file=$(rg -l "image: eirini/${1}")
    new_image_ref="$(docker inspect --format='{{index .RepoDigests 0}}' eirini/opi:patch-me-if-you-can)"
    sed -i '' "s|image: eirini/${1}.*$|image: ${new_image_ref}|g" "$file"
  popd
}

helm_upgrade() {
  pushd "$eirini_release_basedir/helm/cf"
    helm dep update
  popd

  SECRET=$(kubectl get pods --namespace uaa -o jsonpath='{.items[?(.metadata.name=="uaa-0")].spec.containers[?(.name=="uaa")].env[?(.name=="INTERNAL_CA_CERT")].valueFrom.secretKeyRef.name}')
  CA_CERT="$(kubectl get secret $SECRET --namespace uaa -o jsonpath="{.data['internal-ca-cert']}" | base64 --decode -)"

  CLUSTER=$(kubectl config current-context | cut -d / -f 1)
  SECRET_NAME="$(kubectl get secrets | grep "$CLUSTER" | cut -d ' ' -f 1)"
  BITS_TLS_CRT="$(kubectl get secret "$SECRET_NAME" --namespace default -o jsonpath="{.data['tls\.crt']}" | base64 --decode -)"
  BITS_TLS_KEY="$(kubectl get secret "$SECRET_NAME" --namespace default -o jsonpath="{.data['tls\.key']}" | base64 --decode -)"

  pushd "$eirini_release_basedir/helm"
    helm upgrade --install scf ./cf \
      --namespace scf \
      --values "$HOME/workspace/eirini-private-config/environments/kube-clusters/$CLUSTER/scf-config-values.yaml" \
      --set "secrets.UAA_CA_CERT=${CA_CERT}" \
      --set "bits.secrets.BITS_TLS_KEY=${BITS_TLS_KEY}" \
      --set "bits.secrets.BITS_TLS_CRT=${BITS_TLS_CRT}"
  popd
}

for component in "$@"; do
  echo "--- Patching component $component ---"
  docker_build "$component"
  docker_push "$component"
done

update_helm_chart "$component"
helm_upgrade
