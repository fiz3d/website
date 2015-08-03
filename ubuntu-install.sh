#!/bin/bash

sudo apt-get install git

# Save current working directory.
pwd=$(pwd)

# Download and Go and install it to /root/go
tmpdir="/tmp/fiz3d-org-install"
mkdir -p $tmpdir && cd $tmpdir
wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
tar -xzf ./go1.4.2.linux-amd64.tar.gz -C /root

# Cleanup and go back to the previous working directory.
rm -rf $tmpdir
cd $pwd
