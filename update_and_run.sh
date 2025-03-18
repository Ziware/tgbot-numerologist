#!/bin/bash

REPO_URL="https://github.com/Ziware/tgbot-numerologist.git"
REPO_NAME="tgbot-numerologist"
REPO_BRANCH="main"
DOCKER_COMPOSE_FILE="docker-compose.yml"
DOCKER_DIR="src"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

command -v git >/dev/null 2>&1 || error "Git is not installed. Install it using apt, yum or brew."
command -v docker >/dev/null 2>&1 || error "Docker is not installed. Visit https://docs.docker.com/get-docker/"
command -v docker-compose >/dev/null 2>&1 || error "Docker Compose is not installed. Please install it."

if [ -d "$REPO_NAME" ]; then
    log "Updating existing repository $REPO_NAME..."
    cd "$REPO_NAME" || error "Failed to change to directory $REPO_NAME"
    
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    
    git fetch || error "Failed to fetch updates from repository"
    
    if [ "$CURRENT_BRANCH" != "$REPO_BRANCH" ]; then
        log "Switching to branch $REPO_BRANCH..."
        git checkout "$REPO_BRANCH" || error "Failed to switch to branch $REPO_BRANCH"
    fi
    
    git pull || error "Failed to execute git pull"
else
    log "Cloning repository $REPO_URL..."
    git clone --branch "$REPO_BRANCH" "$REPO_URL" "$REPO_NAME" || error "Failed to clone repository"
    
    cd "$REPO_NAME" || error "Failed to change to directory $REPO_NAME"
fi

if [ -n "$DOCKER_DIR" ]; then
    log "Changing to directory $DOCKER_DIR..."
    cd "$DOCKER_DIR" || error "Failed to change to directory $DOCKER_DIR"
fi

log "Starting Docker..."

if [ -f "$DOCKER_COMPOSE_FILE" ] && [ "$HAS_DOCKER_COMPOSE" = "1" ]; then
    log "Stopping existing containers..."
    docker-compose -f "$DOCKER_COMPOSE_FILE" down

    log "Building and starting with Docker Compose..."
    docker-compose -f "$DOCKER_COMPOSE_FILE" up -d --build || error "Error when starting Docker Compose"
    
    log "Docker Compose successfully started!"
else
    error "docker-compose.yml not found. Check the repository."
fi

log "Script successfully completed!"
