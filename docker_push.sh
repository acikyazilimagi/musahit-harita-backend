#!/usr/bin/env bash

# Exit on error
set -e

# Build
AWS_PROFILE=aya-secim aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin 708115278474.dkr.ecr.eu-central-1.amazonaws.com

# Build
docker build -t secim .
docker tag secim:latest 708115278474.dkr.ecr.eu-central-1.amazonaws.com/secim:latest

# Push
docker push 708115278474.dkr.ecr.eu-central-1.amazonaws.com/secim:latest