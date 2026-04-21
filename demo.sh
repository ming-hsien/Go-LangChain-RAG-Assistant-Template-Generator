#!/bin/bash

# demo.sh - Manage Demo Features
# Usage: ./demo.sh [install|clean]

ACTION=${1:-install}
PROJECT_ROOT=$(pwd)
DEMO_SOURCE_DIR="$PROJECT_ROOT/demo/tools"
DEMO_TARGET_DIR="$PROJECT_ROOT/internal/tools/demo"
MAIN_FILE="$PROJECT_ROOT/cmd/app/main.go"
DEMO_IMPORT='_ "github.com/ming-hsien/lang-chain-template/internal/tools/demo"'

if [ "$ACTION" == "install" ]; then
    echo "Installing Demo Tools..."

    mkdir -p "$DEMO_TARGET_DIR"

    if [ -d "$DEMO_SOURCE_DIR" ]; then
        cp "$DEMO_SOURCE_DIR/"* "$DEMO_TARGET_DIR/"
        echo "Demo files copied to internal/tools/demo"
    else
        echo "Error: Demo source directory $DEMO_SOURCE_DIR not found."
        exit 1
    fi

    if ! grep -q "internal/tools/demo" "$MAIN_FILE"; then
        sed -i '/import (/a \	'"$DEMO_IMPORT" "$MAIN_FILE"
        echo "Demo import activated in $MAIN_FILE"
    fi
    
    echo "Demo installation complete! Run 'docker compose up --build' to see it in action."

elif [ "$ACTION" == "clean" ]; then
    echo "Cleaning up Demo Tools..."
    if [ -d "$DEMO_TARGET_DIR" ]; then
        rm -rf "$DEMO_TARGET_DIR"
        echo "Removed internal/tools/demo"
    fi
    
    if grep -q "internal/tools/demo" "$MAIN_FILE"; then
        sed -i "/internal\/tools\/demo/d" "$MAIN_FILE"
        echo "Demo import removed from $MAIN_FILE"
    fi
    
    rm -f documents/*.txt
    echo "Project is now clean."
else
    echo "Usage: ./demo.sh [install|clean]"
    exit 1
fi
