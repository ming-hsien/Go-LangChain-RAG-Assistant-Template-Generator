#!/bin/bash

red=$(tput setaf 1)
green=$(tput setaf 2)
blue=$(tput setaf 4)
reset=$(tput sgr0)

# Configuration
OLD_MODULE="github.com/ming-hsien/lang-chain-template"

if [ -z "$1" ]; then
    echo "Usage: ./init.sh <new_module_name> [destination_directory]"
    echo "Example: ./init.sh github.com/my-org/my-project ./my-new-app"
    exit 1
fi

NEW_MODULE=$1
DEST_DIR=$2

if [ -n "$DEST_DIR" ]; then
    if [ -d "$DEST_DIR" ]; then
        echo -e "${red}Error${reset}: Destination directory '$DEST_DIR' already exists."
        echo "Please remove it or choose a different name to avoid overwriting."
        exit 1
    fi

    echo "Creating destination directory: $DEST_DIR..."
    mkdir -p "$DEST_DIR"
    
    echo "Copying template files (excluding .git and init.sh)..."
    cp -r . "$DEST_DIR"
    rm -rf "$DEST_DIR/.git"
    rm "$DEST_DIR/init.sh"
    
    cd "$DEST_DIR" || exit
fi

echo "Initializing project with new module name: $NEW_MODULE..."

# 1. Update go.mod
if [ -f "go.mod" ]; then
    echo "Updating go.mod..."
    go mod edit -module "$NEW_MODULE"
else
    echo "go.mod not found. Skipping mod edit."
fi

echo "Replacing import paths in .go files..."
grep -rl "$OLD_MODULE" . | xargs -r sed -i "s|$OLD_MODULE|$NEW_MODULE|g"

if [ -f "README.md" ]; then
    echo "Updating README.md references..."
    echo "# $NEW_MODULE" > README.md
fi

echo "Running go mod tidy..."
go mod tidy

echo "======== ${green}Project initialized successfully!${reset} ========"
if [ -n "$DEST_DIR" ]; then
    echo "Location: ${blue}$(pwd)${reset}"
fi
echo "Next steps:"
echo "1. cd ${blue}$(pwd)${reset}"
echo "2. Configure your .env file"
echo "3. Run 'docker compose up -d --build' to start"
echo "==================================================="
