#!/bin/bash -xe

# This script should be able to execute KubeMacPool
# functional tests against Kubernetes cluster with
# CNAO built with latest changes, on any
# environment with basic dependencies listed in
# check-patch.packages installed and docker running.
#
# yum -y install automation/check-patch.packages
# automation/check-patch.e2e-kubemacpool-functests.sh

teardown() {
    make cluster-down
}

main() {
    export KUBEVIRT_PROVIDER='k8s-1.17'

    source automation/check-patch.setup.sh

    # Spin up Kubernetes cluster
    cd ${TMP_PROJECT_PATH}
    make cluster-down
    make cluster-up
    trap teardown EXIT SIGINT SIGTERM SIGSTOP

    # Deploy CNAO with latest changes
    make cluster-operator-push
    make cluster-operator-install

    # Deploy kubevirt
    ./hack/deploy-kubevirt.sh

    # Export kubecinfig absolute path
    export KUBECONFIG=$(_kubevirtci/cluster-up/kubeconfig.sh)

    # Clone KubeMacPool repo
    KUBEMACPOOL_PROJECT_PATH=${GOPATH}/src/github.com/k8snetworkplumbingwg/kubemacpool
    KUBEMACPOOL_URL=$(cat components.yaml | shyaml get-value components.kubemacpool.url)
    KUBEMACPOOL_COMMIT=$(cat components.yaml | shyaml get-value components.kubemacpool.commit)
    mkdir -p $KUBEMACPOOL_PROJECT_PATH
    git clone $KUBEMACPOOL_URL $KUBEMACPOOL_PROJECT_PATH
    cd $KUBEMACPOOL_PROJECT_PATH
    git fetch origin
    git reset --hard
    git checkout ${KUBEMACPOOL_COMMIT}

    # Deploy NetworkAddonsConfig CR
    export KUBEVIRT_PROVIDER=external
    cat <<EOF > cr.yaml
apiVersion: networkaddonsoperator.network.kubevirt.io/v1alpha1
kind: NetworkAddonsConfig
metadata:
  name: cluster
spec:
  imagePullPolicy: Always
  kubeMacPool: {}
  linuxBridge: {}
  multus: {}
EOF
    ./cluster/kubectl.sh apply -f cr.yaml
    ./cluster/kubectl.sh wait networkaddonsconfig cluster --for condition=Available --timeout=800s

    # Run KubeMacPool functional tests
    make functest
}

[[ "${BASH_SOURCE[0]}" == "$0" ]] && main "$@"
