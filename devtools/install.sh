#!/bin/bash
set -euo pipefail

# Script version
SCRIPT_VERSION="1.0.0"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Required versions
REQUIRED_GO_VERSION="1.24.0"
REQUIRED_NODE_VERSION="20.0.0"
REQUIRED_PNPM_VERSION="10.8.0"
GOLANGCI_LINT_VERSION="1.59.1"

# Platform detection
detect_platform() {
    case "$(uname -s)" in
        Darwin*)    PLATFORM="mac" ;;
        Linux*)     
            if [ -f /etc/os-release ]; then
                . /etc/os-release
                if [[ "$ID" == "ubuntu" ]]; then
                    PLATFORM="ubuntu"
                else
                    PLATFORM="linux"
                fi
            else
                PLATFORM="linux"
            fi
            ;;
        CYGWIN*|MINGW*|MSYS*)    PLATFORM="win" ;;
        *)          
            echo -e "${RED}Unsupported platform: $(uname -s)${NC}"
            exit 1
            ;;
    esac
}

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Version comparison function
version_ge() {
    local version1=$1
    local version2=$2
    printf '%s\n%s\n' "$version2" "$version1" | sort -V -C
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download with checksum verification
secure_download() {
    local url=$1
    local output=$2
    local checksum_url=$3
    local expected_checksum=$4
    
    log_info "Downloading from $url..."
    
    # Create temp directory
    local tmp_dir="$(mktemp -d)"
    trap 'rm -rf "$tmp_dir"' EXIT
    
    # Download file
    if ! curl -fsSL "$url" -o "$tmp_dir/$(basename "$output")"; then
        log_error "Failed to download $url"
        return 1
    fi
    
    # Verify checksum if provided
    if [ -n "$checksum_url" ] || [ -n "$expected_checksum" ]; then
        log_info "Verifying checksum..."
        local actual_checksum
        
        if [ -n "$checksum_url" ]; then
            # Download checksum file
            curl -fsSL "$checksum_url" -o "$tmp_dir/checksum.txt" || {
                log_error "Failed to download checksum"
                return 1
            }
            expected_checksum="$(awk '{print $1}' "$tmp_dir/checksum.txt")"
        fi
        
        # Calculate actual checksum
        if command_exists sha256sum; then
            actual_checksum="$(sha256sum "$tmp_dir/$(basename "$output")" | awk '{print $1}')"
        elif command_exists shasum; then
            actual_checksum="$(shasum -a 256 "$tmp_dir/$(basename "$output")" | awk '{print $1}')"
        else
            log_warn "Cannot verify checksum - no sha256sum or shasum available"
            mv "$tmp_dir/$(basename "$output")" "$output"
            return 0
        fi
        
        if [ "$actual_checksum" != "$expected_checksum" ]; then
            log_error "Checksum verification failed!"
            log_error "Expected: $expected_checksum"
            log_error "Actual: $actual_checksum"
            return 1
        fi
        
        log_success "Checksum verified"
    fi
    
    # Move to final location
    mv "$tmp_dir/$(basename "$output")" "$output"
    return 0
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    local prerequisites_met=true
    
    # Check curl
    if ! command_exists curl; then
        log_error "curl is required but not installed"
        prerequisites_met=false
    fi
    
    # Check git
    if ! command_exists git; then
        log_error "git is required but not installed"
        prerequisites_met=false
    fi
    
    # Check Go version
    if command_exists go; then
        local go_version="$(go version | awk '{print $3}' | sed 's/go//')"
        if ! version_ge "$go_version" "$REQUIRED_GO_VERSION"; then
            log_error "Go $REQUIRED_GO_VERSION or higher is required (found: $go_version)"
            prerequisites_met=false
        else
            log_success "Go $go_version found"
        fi
    else
        log_error "Go $REQUIRED_GO_VERSION or higher is required but not installed"
        log_info "Please install Go from https://golang.org/dl/"
        prerequisites_met=false
    fi
    
    # Check Node.js version
    if command_exists node; then
        local node_version="$(node -v | sed 's/v//')"
        if ! version_ge "$node_version" "$REQUIRED_NODE_VERSION"; then
            log_error "Node.js $REQUIRED_NODE_VERSION or higher is required (found: $node_version)"
            prerequisites_met=false
        else
            log_success "Node.js $node_version found"
        fi
    else
        log_error "Node.js $REQUIRED_NODE_VERSION or higher is required but not installed"
        log_info "Please install Node.js from https://nodejs.org/"
        prerequisites_met=false
    fi
    
    # Check npm
    if ! command_exists npm; then
        log_error "npm is required but not installed"
        prerequisites_met=false
    fi
    
    if [ "$prerequisites_met" = false ]; then
        log_error "Prerequisites check failed. Please install missing dependencies."
        exit 1
    fi
    
    log_success "All prerequisites met"
}

# Check and setup PATH
setup_path() {
    local go_bin="$HOME/go/bin"
    local path_updated=false
    
    # Check if go/bin is in PATH
    if [[ ":$PATH:" != *":$go_bin:"* ]]; then
        log_warn "$go_bin is not in PATH"
        
        # Detect shell and update appropriate config file
        local shell_config=""
        if [ -n "$BASH_VERSION" ]; then
            shell_config="$HOME/.bashrc"
        elif [ -n "$ZSH_VERSION" ]; then
            shell_config="$HOME/.zshrc"
        else
            shell_config="$HOME/.profile"
        fi
        
        # Add to PATH
        echo "" >> "$shell_config"
        echo "# Added by Sourcetool install script" >> "$shell_config"
        echo "export PATH=\"\$PATH:$go_bin\"" >> "$shell_config"
        
        log_info "Added $go_bin to PATH in $shell_config"
        log_warn "Please run 'source $shell_config' or restart your terminal"
        
        # Export for current session
        export PATH="${PATH}:${go_bin}"
        path_updated=true
    fi
    
    return 0
}

# Install pnpm
install_pnpm() {
    if command_exists pnpm; then
        local current_version="$(pnpm -v)"
        if version_ge "$current_version" "$REQUIRED_PNPM_VERSION"; then
            log_success "pnpm $current_version already installed"
            return 0
        else
            log_info "Upgrading pnpm from $current_version to $REQUIRED_PNPM_VERSION"
        fi
    fi
    
    log_info "Installing pnpm v$REQUIRED_PNPM_VERSION..."
    if npm install -g "pnpm@$REQUIRED_PNPM_VERSION"; then
        log_success "pnpm installed successfully"
    else
        log_error "Failed to install pnpm"
        return 1
    fi
}

# Install golangci-lint
install_golangci_lint() {
    if command_exists golangci-lint; then
        log_success "golangci-lint already installed"
        return 0
    fi
    
    log_info "Installing golangci-lint v$GOLANGCI_LINT_VERSION..."
    
    case "$PLATFORM" in
        mac)
            if command_exists brew; then
                brew install golangci-lint
            else
                # Manual install
                curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
                    sh -s -- -b "$HOME/go/bin" "v$GOLANGCI_LINT_VERSION"
            fi
            ;;
        ubuntu|linux)
            curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
                sh -s -- -b "$HOME/go/bin" "v$GOLANGCI_LINT_VERSION"
            ;;
        win)
            if command_exists scoop; then
                scoop install golangci-lint
            elif command_exists choco; then
                choco install golangci-lint
            else
                log_warn "Please install golangci-lint manually from:"
                log_warn "https://github.com/golangci/golangci-lint/releases/tag/v$GOLANGCI_LINT_VERSION"
                return 1
            fi
            ;;
    esac
    
    if command_exists golangci-lint; then
        log_success "golangci-lint installed successfully"
    else
        log_error "Failed to install golangci-lint"
        return 1
    fi
}

# Install gofumpt
install_gofumpt() {
    if command_exists gofumpt; then
        log_success "gofumpt already installed"
        return 0
    fi
    
    log_info "Installing gofumpt..."
    if go install mvdan.cc/gofumpt@latest; then
        log_success "gofumpt installed successfully"
    else
        log_error "Failed to install gofumpt"
        return 1
    fi
}

# Install buf
install_buf() {
    if command_exists buf; then
        log_success "buf already installed"
        return 0
    fi
    
    log_info "Installing buf..."
    
    case "$PLATFORM" in
        mac)
            if command_exists brew; then
                brew install bufbuild/buf/buf
            else
                # Manual install
                local buf_bin="$HOME/go/bin/buf"
                secure_download \
                    "https://github.com/bufbuild/buf/releases/latest/download/buf-Darwin-x86_64" \
                    "$buf_bin" \
                    "" ""
                chmod +x "$buf_bin"
            fi
            ;;
        ubuntu|linux)
            local buf_bin="$HOME/go/bin/buf"
            local arch="$(uname -m)"
            case "$arch" in
                x86_64) arch="x86_64" ;;
                aarch64) arch="aarch64" ;;
                *) 
                    log_error "Unsupported architecture: $arch"
                    return 1
                    ;;
            esac
            
            secure_download \
                "https://github.com/bufbuild/buf/releases/latest/download/buf-Linux-${arch}" \
                "$buf_bin" \
                "" ""
            chmod +x "$buf_bin"
            ;;
        win)
            if command_exists scoop; then
                scoop install buf
            elif command_exists choco; then
                choco install buf
            else
                log_warn "Please install buf manually from:"
                log_warn "https://github.com/bufbuild/buf/releases"
                return 1
            fi
            ;;
    esac
    
    if command_exists buf; then
        log_success "buf installed successfully"
    else
        log_error "Failed to install buf"
        return 1
    fi
}

# Install Docker
install_docker() {
    if command_exists docker; then
        log_success "Docker already installed"
        return 0
    fi
    
    log_info "Installing Docker..."
    
    case "$PLATFORM" in
        mac)
            if command_exists brew; then
                brew install --cask docker
                log_info "Please start Docker Desktop from Applications"
            else
                log_warn "Please install Docker Desktop manually from:"
                log_warn "https://www.docker.com/products/docker-desktop/"
                return 1
            fi
            ;;
        ubuntu)
            # Docker official installation script
            log_info "Installing Docker using official script..."
            curl -fsSL https://get.docker.com | sudo sh
            
            # Add user to docker group
            sudo usermod -aG docker "${USER}"
            log_warn "You need to log out and back in for docker group membership to take effect"
            ;;
        linux)
            log_warn "Please install Docker for your distribution from:"
            log_warn "https://docs.docker.com/engine/install/"
            return 1
            ;;
        win)
            log_warn "Please install Docker Desktop manually from:"
            log_warn "https://www.docker.com/products/docker-desktop/"
            return 1
            ;;
    esac
    
    if command_exists docker; then
        log_success "Docker installed successfully"
    else
        log_warn "Docker installation requires manual steps or system restart"
    fi
}

# Main installation flow
main() {
    echo "==================================="
    echo "Sourcetool Development Setup v$SCRIPT_VERSION"
    echo "==================================="
    echo
    
    # Detect platform
    detect_platform
    log_info "Detected platform: $PLATFORM"
    
    # Parse arguments
    if [ $# -eq 1 ]; then
        case "$1" in
            mac|ubuntu|linux|win)
                PLATFORM="$1"
                log_info "Platform override: $PLATFORM"
                ;;
            --help|-h)
                echo "Usage: $0 [mac|ubuntu|linux|win]"
                echo
                echo "Installs development dependencies for Sourcetool"
                echo "Platform is auto-detected if not specified"
                exit 0
                ;;
            *)
                log_error "Invalid platform: $1"
                echo "Usage: $0 [mac|ubuntu|linux|win]"
                exit 1
                ;;
        esac
    fi
    
    # Check prerequisites
    check_prerequisites
    
    # Setup PATH
    setup_path
    
    # Track installation results
    local failed_installs=()
    
    # Install tools
    echo
    log_info "Installing development tools..."
    
    install_pnpm || failed_installs+=("pnpm")
    install_golangci_lint || failed_installs+=("golangci-lint")
    install_gofumpt || failed_installs+=("gofumpt")
    install_buf || failed_installs+=("buf")
    install_docker || failed_installs+=("docker")
    
    # Summary
    echo
    echo "==================================="
    echo "Installation Summary"
    echo "==================================="
    
    if [ ${#failed_installs[@]} -eq 0 ]; then
        log_success "All tools installed successfully!"
        
        # Next steps
        echo
        echo "Next steps:"
        echo "1. If PATH was updated, run: source ~/.bashrc (or ~/.zshrc)"
        echo "2. If Docker was installed, ensure Docker daemon is running"
        echo "3. Run 'make dev-ce' to start the development environment"
        echo "4. Run 'pnpm install' in the project root"
    else
        log_warn "Some tools failed to install:"
        for tool in "${failed_installs[@]}"; do
            echo "  - $tool"
        done
        echo
        echo "Please install these tools manually and run the script again"
        exit 1
    fi
}

# Run main function
main "$@"