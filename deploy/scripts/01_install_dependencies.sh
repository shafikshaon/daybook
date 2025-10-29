#!/bin/bash
set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Load configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/../config/deploy.conf"

if [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
else
    log_error "Configuration file not found: $CONFIG_FILE"
    exit 1
fi

log_info "Installing system dependencies..."

# Update system packages
log_info "Updating system packages..."
sudo apt-get update
sudo apt-get upgrade -y

# Install basic dependencies
log_info "Installing basic tools..."
sudo apt-get install -y \
    curl \
    wget \
    git \
    build-essential \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release \
    ufw \
    fail2ban \
    htop \
    vim \
    unzip

# Install Go
log_info "Installing Go ${GO_VERSION}..."
if command -v go &> /dev/null; then
    CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_warn "Go $CURRENT_GO_VERSION is already installed"
    read -p "Do you want to reinstall Go ${GO_VERSION}? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Skipping Go installation"
    else
        sudo rm -rf /usr/local/go
        wget "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
        sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
        rm "go${GO_VERSION}.linux-amd64.tar.gz"
    fi
else
    wget "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
    rm "go${GO_VERSION}.linux-amd64.tar.gz"
fi

# Add Go to PATH if not already added
if ! grep -q "/usr/local/go/bin" ~/.profile; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.profile
    source ~/.profile
fi

export PATH=$PATH:/usr/local/go/bin
log_info "Go installed: $(go version)"

# Install Node.js
log_info "Installing Node.js ${NODE_VERSION}..."
if command -v node &> /dev/null; then
    CURRENT_NODE_VERSION=$(node -v | sed 's/v//' | cut -d. -f1)
    log_warn "Node.js v$CURRENT_NODE_VERSION is already installed"
else
    curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | sudo -E bash -
    sudo apt-get install -y nodejs
    log_info "Node.js installed: $(node -v)"
    log_info "npm installed: $(npm -v)"
fi

log_info "System dependencies installed successfully!"

# Display versions
log_info "Installed versions:"
echo "  - Go: $(go version | awk '{print $3}')"
echo "  - Node.js: $(node -v)"
echo "  - npm: $(npm -v)"

exit 0
