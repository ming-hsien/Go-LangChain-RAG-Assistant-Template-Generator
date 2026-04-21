# Generic RAG AI Assistant Template

[![Go Version](https://img.shields.io/badge/Go-1.26.2-blue.svg)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Supported-blue.svg)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go-based **Retrieval-Augmented Generation (RAG)** template built with **LangChainGo**, **GitHub Models**, and **Qdrant**. Designed for rapid deployment of custom AI assistants with internal knowledge retrieval.

## Key Features

- **High-Performance RAG**: Powered by `langchaingo` for seamless document chunking, embedding, and retrieval workflows.
- **GitHub Models Integration**: Direct support for GitHub Marketplace models (GPT-4o & text-embedding-3-small) using standard OpenAI-compatible endpoints.
- **Auto-Provisioning Vector Store**: Automatically checks and initializes Qdrant collections (CMS_documents) on startup.
- **Modern Web Interface**: Built-in glassmorphism-style UI using TailwindCSS and vanilla JS.
- **One-Click Deployment**: Dockerized multi-stage builds with Docker Compose for rapid environment setup.

---

## Quick Start

### 1. Prerequisites
- Get a Personal Access Token (PAT) from [GitHub Models](https://github.com/marketplace/models).
- Ensure Docker and Docker Compose are installed.

### 2. Project Initialization (Template Generator)
To generate a new project from this template:
```bash
./init.sh github.com/your-username/your-project-name [destination_directory]
```
- `<new_module_name>`: The Go module path for your new project.
- `[destination_directory]` (Optional): Where to create the new project. If omitted, it will initialize in the current directory.

This script will copy the files (excluding `.git`), update `go.mod`, and refactor all internal import paths automatically.

### 3. Configuration
Copy the environment template and fill in your details:
```bash
cp .env.example .env
```
Edit `.env`:
```env
GITHUB_TOKEN=your_github_pat_here
APP_NAME=My AI Assistant
COLLECTION_NAME=my_documents
SYSTEM_PROMPT=You are a helpful assistant...
```

### 3. Deploy with Docker Compose (Recommended)
```bash
docker-compose up -d --build
```
Once started, visit: `http://localhost:8080`

### 4. Local Development
Ensure you have Go 1.26.2+ installed and Qdrant is running:
```bash
go run cmd/app/main.go
```

---

## Project Structure

```text
.
├── cmd/
│   └── app/           # Main entry point (main.go)
├── internal/
│   ├── api/           # Gin Web Server and static UI assets
│   ├── config/        # Environment variable and config loader
│   ├── promptmgr/     # (Placeholder) Prompt management logic
│   └── rag/           # Core RAG logic (Vector Search, Embedding, LLM)
├── documents/         # Place your .txt files here for indexing
├── prompts/           # (Placeholder) Structured prompt templates
├── Dockerfile         # Production-ready multi-stage Docker build
└── docker-compose.yml # Service orchestration (App Core + Qdrant)
```

---

## Operations

### Knowledge Indexing
Place your network configuration or guide files in the `documents/` folder, then click **"Re-Index Knowledge"** in the Web UI. The system will:
1. Parse all documents.
2. Chunk text using recursive splitting.
3. Generate embeddings via OpenAI.
4. Upsert vectors into Qdrant.

### API Endpoints
- `POST /v1/query`: Submit a question to get a RAG-powered answer.
- `POST /v1/index`: Trigger a background re-indexing task.
- `GET /v1/tools`: (Mock) Retrieve availability list of tools for AI Agents.

---

## License
This project is licensed under the [MIT License](LICENSE).
