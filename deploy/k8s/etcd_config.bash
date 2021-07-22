#!/usr/bin/env bash


MINIKUBE_IP=$(minikube ip)
MINIKUBE_PORT=$(kubectl get services -o jsonpath="{.spec.ports[].nodePort}" etcd-client -n etcd)

cat - |\
sed 's/\$MINIKUBE_IP'"/$MINIKUBE_IP/g" |\
sed 's/\$MINIKUBE_PORT'"/$MINIKUBE_PORT/g"