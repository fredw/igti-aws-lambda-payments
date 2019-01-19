#!/usr/bin/env bash

# Install AWS SDK
pip install --user awscli
export PATH=$PATH:$HOME/.local/bin

# AWS Login
eval $(aws ecr get-login)

# Build and upload the function on AWS
make -C $TRAVIS_BUILD_DIR build
make -C $TRAVIS_BUILD_DIR update

# Update the handler name to the right name
#aws lambda update-function-configuration --function-name=payments --handler=main
