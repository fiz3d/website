#!/bin/bash

sudo apt-get install git

# Save current working directory.
export PWD=$(pwd)

# Download and Go and install it to /root/go
export TMPDIR="/tmp/fiz3d-org-install"
mkdir -p $TMPDIR && cd $TMPDIR
wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
mkdir -p /root/go
tar -xzf ./go1.4.2.linux-amd64.tar.gz -C /root/go

# Cleanup and go back to the previous working directory.
rm -rf $TMPDIR
cd $PWD
