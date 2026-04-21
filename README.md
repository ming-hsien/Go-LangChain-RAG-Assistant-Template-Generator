# Go-LangChain RAG Assistant Template

[![Go Version](https://img.shields.io/badge/Go-1.26.2-blue.svg)](https://go.dev/)
[![LangChainGo](https://img.shields.io/badge/Powered%20by-LangChainGo-green.svg)](https://github.com/tmc/langchaingo)
[![Docker](https://img.shields.io/badge/Docker-Supported-blue.svg)](https://www.docker.com/)
[![Architecture](https://img.shields.io/badge/Architecture-SBI--Processor--RAG-orange.svg)](#project-structure)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go-based **Retrieval-Augmented Generation (RAG)** assistant template, designed with **LangChainGo** and following **SBI (Service-Based Interface)** architecture principles.

## Architecture Concept: LangChain & SBI

This template centers around **LangChain (langchaingo)**, implementing AI orchestration through a modular, service-based approach:
- **LangChain Orchestration**: The **Processor** acts as a LangChain-style "Chain", orchestrating models, vector stores, and logic to handle semantic queries.
- **SBI (Service-Based Interface)**: Following 5G core network design patterns, the SBI exposes these LangChain-powered capabilities as robust, machine-readable services.
- **Decoupled Components**: Leveraging LangChain's modularity, the system easily separates Data Indexing from Inference logic.

---

## Key Features

- **Decoupled SBI Layer**: Professional project structure with clear separation between transmission, logic, and data.
- **High-Performance RAG**: Powered by `langchaingo` for seamless document chunking and retrieval workflows.
- **GitHub Models Integration**: Direct support for GPT-4o & text-embedding-3-small using OpenAI-compatible endpoints.
- **Independent Web Interface**: A modern glassmorphism UI located in a dedicated `/web` directory, decoupled from the API logic.
- **One-Click Deployment**: Dockerized multi-stage builds with Docker Compose.

---

## Quick Start

### 1. Prerequisites
- Get a Personal Access Token (PAT) from [GitHub Models](https://github.com/marketplace/models).
- Ensure Docker and Docker Compose are installed.

### 2. Project Initialization
Generate a new project from this template:
```bash
./gen.sh github.com/your-username/your-project-name [destination_directory]
```
This script handles module renaming, import path refactoring, and directory preparation automatically.

### 3. Configuration
Copy the environment template and fill in your details:
```bash
cp .env.example .env
```
Key parameters in `.env`:
- `GITHUB_TOKEN`: Your GitHub Personal Access Token for Model API.
- `APP_NAME`: Your AI assistant's name (displays on UI).
- `COLLECTION_NAME`: Qdrant collection name for your knowledge base.
- `SYSTEM_PROMPT`: Directs the AI's behavior and persona.

### 4. Deploy with Docker (Recommended)
```bash
docker-compose up -d --build
```
Access the UI at: `http://localhost:8080/ui/`

### 5. Local Development
Ensure you have Go 1.26.2+ installed and Qdrant is running:
```bash
# Start Qdrant container if not running
docker run -d -p 6333:6333 qdrant/qdrant

# Run the application
go run cmd/app/main.go
```

## Project Structure

```text
.
├── cmd/
│   └── app/           # Main entry point (Wiring SBI + Processor + RAG)
├── internal/
│   ├── sbi/           # Service-Based Interface (API Routing & Handlers)
│   │   ├── processor/ # Internal Business Logic & Agentic reasoning
│   │   └── server.go  # Web Server & Static Asset Mounting
│   ├── tools/         # Dynamic Tool Framework (Add tools here)
│   ├── rag/           # Low-level RAG logic (Vector DB, Embedding, LLM)
│   └── config/        # Configuration management
├── web/               # Decoupled Web Frontend (HTML, JS, Assets)
├── documents/         # Local knowledge base (.txt files)
├── Dockerfile         # Multi-stage production build
└── docker-compose.yml # Service orchestration
```

---

## Quick Demo (Optional AI Agent features)
This template includes a pre-built demo showcasing **Agentic Tool Calling** (System monitoring & File generation).
- **To Install Demo**: `./demo.sh install` (This activates the demo files and logic).
- **To Clean Demo**: `./demo.sh clean` (This completely removes demo files and reverts the project to a clean template state).

### AI Agent Capabilities (Tools)
When the **Quick Demo** is installed, the AI Agent gains the following multi-step reasoning capabilities:
- **System Monitoring**: Access to real-time host memory and CPU load stats.
- **Document Generation**: Ability to write new summary files directly into the `documents/` directory.

---

## API Endpoints (SBI)
- `POST /v1/query`: Submit a question to the Processor (supports RAG + Agent tools).
- `POST /v1/index`: Trigger background re-indexing of the knowledge base.
- `GET /v1/config`: Retrieve dynamic UI branding and system status.
- `GET /v1/tools`: Retrieve the list of currently registered AI Agent tools.

---

## License
This project is licensed under the [MIT License](LICENSE).
