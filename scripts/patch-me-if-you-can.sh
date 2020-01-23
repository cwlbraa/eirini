#!/bin/bash

set -xeuo pipefail

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
eirini_basedir=$(realpath "$script_dir/..")
eirini_release_basedir=$(realpath "$script_dir/../../eirini-release")

docker_build() {
  echo "Building docker image for $1"
  pushd "$eirini_basedir"
    docker build . -f "$eirini_basedir/docker/$component/Dockerfile" \
      --build-arg GIT_SHA=dev \
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
    sed -i '' "s|image: eirini/${1}.*$|image: eirini/${1}:patch-me-if-you-can|g" "$file"
  popd

}

helm_upgrade() {
  pushd "$eirini_release_basedir/helm/cf"
    helm dep update
  popd

  # pushd "$eirini_release_basedir/helm"
  #   helm upgrade --install scf ./cf \
  #     --namespace scf \
  #     --values ~/workspace/eirini-private-config/environments/kube-clusters/veliko-tarnovo/scf-config-values.yaml \
  #     --set "secrets.UAA_CA_CERT=${CA_CERT}" \
  #     --set "bits.secrets.BITS_TLS_KEY=${BITS_TLS_KEY}" \
  #     --set "bits.secrets.BITS_TLS_CRT=${BITS_TLS_CRT}"
  # popd
}

for component in "$@"; do
  echo "--- Patching component $component ---"
  # docker_build "$component"
  # docker_push "$component"
  update_helm_chart "$component"
done

helm_upgrade
