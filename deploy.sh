#!/bin/bash

# Copyright (C) 2023 Jesus Rodriguez Miranda (jesus@wallib.com)

set -Eeuo pipefail
set -o xtrace

export NAMESPACE=
export SERVICE_NAME=
export APP_NAME=
export SERVICE_PORT=8080
export SERVICE_TARGET_PORT=8080
envsubst < k8s.yml > k8s-deploy.yml
k8s-deploy.yml

kubectl apply -f k8s-deploy.yml
