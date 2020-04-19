#!/bin/bash

echo "Creating services..."

kubectl create -f ingester-es-svc.yaml
kubectl create -f ingester-configmap-svc.yaml
kubectl create -f ingester-defaultapp-svc.yaml

sleep 5s

kubectl get svc -o wide

echo "Done"
