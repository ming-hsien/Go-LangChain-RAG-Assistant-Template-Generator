#!/bin/bash

# Configuration
OLD_MODULE="github.com/ming-hsien/lang-chain-template"

# Check if a new module name was provided
if [ -z "$1" ]; then
    echo "Usage: ./init.sh <new_module_name>"
    echo "Example: ./init.sh github.com/my-org/my-project"
    exit 1
fi

NEW_MODULE=$1

echo "Initializing project with new module name: $NEW_MODULE..."

# 1. Update go.mod
if [ -f "go.mod" ]; then
    echo "Updating go.mod..."
    go mod edit -module "$NEW_MODULE"
else
    echo "go.mod not found. Skipping mod edit."
fi

# 2. Update all .go files recursively
echo "Replacing import paths in .go files..."
find . -type f -name "*.go" -print0 | xargs -0 sed -i "s|$OLD_MODULE|$NEW_MODULE|g"

# 3. Update other relevant files (e.g., README.md references)
if [ -f "README.md" ]; then
    echo "Updating README.md references..."
    sed -i "s|$OLD_MODULE|$NEW_MODULE|g"
fi

# 4. Tidy up the Go module
echo "Running go mod tidy..."
go mod tidy

echo "======== Project initialized successfully! ========"
echo "Next steps:"
echo "1. Configure your .env file"
echo "2. Run 'docker-compose up -d --build' to start"
echo "==================================================="
