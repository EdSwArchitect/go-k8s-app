#!/bin/bash

echo "Creating configmap"

kubectl create configmap ingester-config --from-file=config

kubectl get configmap ingester-config -o yaml

echo "Creating pods..."

kubectl create -f es-pod.yaml
kubectl create -f ingester-configmap-pod.yaml
kubectl create -f ingester-defaultapp-pod.yaml

sleep 5s

kubectl get pod -o wide

echo "Done"
