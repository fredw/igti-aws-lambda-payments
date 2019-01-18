#!/usr/bin/env bash

# Install AWS SDK
pip install --user awscli
export PATH=$PATH:$HOME/.local/bin

# AWS Login
eval $(aws ecr get-login)

# Update the handler name to the right name
aws lambda update-function-configuration --function-name=payments --handler=main
