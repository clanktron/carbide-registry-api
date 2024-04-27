#!/bin/sh
# dependencies: kubectl

usage() {
    echo "Usage: $0 [--connect]"
    echo "Options:"
    echo "  --connect          Connect to the InnoDB Cluster using MySQL Shell"
    echo "  --help             Show this usage message"
}

# NS to deploy carbide-registry-api db and api
NAMESPACE="${NAMESPACE:-carbide-registry-api}"
CREDENTIALS_SECRET_NAME=mypwds
INNODB_CLUSTER_NAME=carbide

connectToCluster() {
    kubectl run --rm -it myshell --image=container-registry.oracle.com/mysql/community-operator -- mysqlsh --sql
}

ensureNamespace() {
if ! kubectl get ns "$NAMESPACE" > /dev/null 2>&1; then
    kubectl create ns "$NAMESPACE"
fi
}

ensureCRDS() {
if ! kubectl get customresourcedefinitions.apiextensions.k8s.io innodbclusters.mysql.oracle.com > /dev/null 2>&1; then
    kubectl apply -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/deploy/deploy-crds.yaml
fi
}

ensureOperator() {
if ! kubectl get deployment mysql-operator -n mysql-operator > /dev/null 2>&1; then
    kubectl apply -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/deploy/deploy-operator.yaml
fi
}

createSecret() {
kubectl create secret generic "$CREDENTIALS_SECRET_NAME" \
        --from-literal=rootUser=root \
        --from-literal=rootHost=% \
        --from-literal=rootPassword="carbidecarbidecarbide" -n "$NAMESPACE"
}

createInnoDBCluster() {
kubectl apply -n "$NAMESPACE" -f - <<EOF
apiVersion: mysql.oracle.com/v2
kind: InnoDBCluster
metadata:
  name: $INNODB_CLUSTER_NAME
spec:
  secretName: $CREDENTIALS_SECRET_NAME
  tlsUseSelfSigned: true
  instances: 3
  router:
    instances: 1
EOF
}

connectFlag=false

while [ $# -gt 0 ]; do
    case "$1" in
        --connect)
            connectFlag=true
            ;;
         --help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
    shift
done

bootstrapOperator() {
    ensureNamespace
    ensureCRDS
    ensureOperator
}

bootstrapCluster() {
    createSecret
    createInnoDBCluster
}

main() {
    bootstrapOperator
    if [ "$connectFlag" = "true" ]; then
        connectToCluster
        exit
    else
        bootstrapCluster
    fi
}

main
