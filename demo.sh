#!/bin/bash

# demo.sh - Manage Demo Features
# Usage: ./demo.sh [install|clean]

ACTION=${1:-install}
PROJECT_ROOT=$(pwd)

# Source directories
SOURCE_TOOLS="$PROJECT_ROOT/demo/tools"
SOURCE_SERVICE="$PROJECT_ROOT/demo/service"
SOURCE_CLIENT="$PROJECT_ROOT/demo/employee_status"
SOURCE_PROMPT="$PROJECT_ROOT/demo/prompt/system.txt"

# Target directories
TARGET_TOOLS="$PROJECT_ROOT/internal/tools/demo"
TARGET_SERVICE="$PROJECT_ROOT/internal/service"
TARGET_CLIENT="$PROJECT_ROOT/internal/employee_status"
TARGET_PROMPT="$PROJECT_ROOT/internal/promptmgr/system.txt"

MAIN_FILE="$PROJECT_ROOT/cmd/app/main.go"
DEMO_IMPORT='_ "github.com/ming-hsien/lang-chain-template/internal/tools/demo"'

if [ "$ACTION" == "install" ]; then
    echo "Installing Demo Features (Tools + Service + Client + Prompt)..."

    # 1. Install Tools
    mkdir -p "$TARGET_TOOLS"
    if [ -d "$SOURCE_TOOLS" ]; then
        cp "$SOURCE_TOOLS/"* "$TARGET_TOOLS/"
        echo " - Tools copied to internal/tools/demo"
    fi

    # 2. Install Service
    mkdir -p "$TARGET_SERVICE"
    if [ -d "$SOURCE_SERVICE" ]; then
        cp "$SOURCE_SERVICE/"* "$TARGET_SERVICE/"
        echo " - Service copied to internal/service"
    fi

    # 3. Install Client
    mkdir -p "$TARGET_CLIENT"
    if [ -d "$SOURCE_CLIENT" ]; then
        cp "$SOURCE_CLIENT/"* "$TARGET_CLIENT/"
        echo " - Client copied to internal/employee_status"
    fi

    # 4. Install Prompt
    if [ -f "$SOURCE_PROMPT" ]; then
        cp "$SOURCE_PROMPT" "$TARGET_PROMPT"
        echo " - Demo system prompt installed."
    fi

    # 5. Activate Import in main.go
    if ! grep -q "internal/tools/demo" "$MAIN_FILE"; then
        sed -i '/import (/a \	'"$DEMO_IMPORT" "$MAIN_FILE"
        echo " - Demo import activated in $MAIN_FILE"
    fi
    
    # 6. Build and Run with Docker
    echo "Starting Docker services..."
    docker compose up --build

    echo "Demo installation complete! Access the UI at: http://localhost:8080/ui/"

elif [ "$ACTION" == "clean" ]; then
    echo "Cleaning up Demo Features..."

    # 0. Shut down Docker services
    echo "Stopping Docker services..."
    docker compose down

    # 1. Remove Tools
    if [ -d "$TARGET_TOOLS" ]; then
        rm -rf "$TARGET_TOOLS"
        echo " - Removed internal/tools/demo"
    fi
    
    # 2. Remove Service (Only specific demo files)
    if [ -f "$TARGET_SERVICE/employee_service.go" ]; then
        rm -f "$TARGET_SERVICE/employee_service.go"
        echo " - Removed internal/service/employee_service.go"
    fi

    # 3. Remove Client
    if [ -d "$TARGET_CLIENT" ]; then
        rm -rf "$TARGET_CLIENT"
        echo " - Removed internal/employee_status"
    fi

    # 4. Remove Import from main.go
    if grep -q "internal/tools/demo" "$MAIN_FILE"; then
        sed -i "/internal\/tools\/demo/d" "$MAIN_FILE"
        echo " - Demo import removed from $MAIN_FILE"
    fi

    # 5. Reset Prompt to generic default
    echo "You are a helpful AI assistant." > "$TARGET_PROMPT"
    echo " - System prompt reset to default."
    
    # 6. Clean documents (Optional: preserve company_rules if you want)
    rm -f documents/*.txt documents/*.md
    echo " - Knowledge base documents cleared."

    echo "Project is now clean."
else
    echo "Usage: ./demo.sh [install|clean]"
    exit 1
fi
