#!/bin/bash
set -e

# Update package lists
sudo apt-get update

# Install Go 1.17 or later
GO_VERSION="1.21.0"
GO_TARBALL="go${GO_VERSION}.linux-amd64.tar.gz"

# Download and install Go
cd /tmp
wget -q "https://golang.org/dl/${GO_TARBALL}"
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "${GO_TARBALL}"

# Add Go to PATH in user's profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
echo 'export GOPATH=$HOME/go' >> $HOME/.profile
echo 'export PATH=$PATH:$GOPATH/bin' >> $HOME/.profile

# Set up Go environment for current session
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Verify Go installation
go version

# Navigate to the workspace directory
cd /mnt/persist/workspace

# Download dependencies
go mod download
go mod tidy

# Verify the module is properly set up
go mod verify