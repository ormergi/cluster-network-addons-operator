#!/bin/bash -xe

KUBEVIRT_VERSION=${KUBEVIRT_VERSION:-v0.29.0}

./cluster/kubectl.sh apply -f https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-operator.yaml

# Ensure the KubeVirt CRD is created
count=0
until ./cluster/kubectl.sh get crd kubevirts.kubevirt.io; do
    ((count++)) && ((count == 30)) && echo "KubeVirt CRD not found" && exit 1
    echo "waiting for KubeVirt CRD"
    sleep 1
done

./cluster/kubectl.sh apply -f https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-cr.yaml

# Ensure the KubeVirt CR is created
count=0
until ./cluster/kubectl.sh -n kubevirt get kv kubevirt; do
    ((count++)) && ((count == 30)) && echo "KubeVirt CR not found" && exit 1
    echo "waiting for KubeVirt CR"
    sleep 1
done

./cluster/kubectl.sh wait -n kubevirt kv kubevirt --for condition=Available --timeout 180s || (echo "KubeVirt not ready in time" && exit 1)

echo "Done"
