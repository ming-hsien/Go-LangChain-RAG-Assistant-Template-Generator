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
- **Agentic Loop**: Supports iterative reasoning (Reasoning-Action-Observation), allowing the AI to call tools, observe results, and refine its answer in multiple turns.
- **SBI (Service-Based Interface)**: Following 5G core network design patterns, the SBI exposes these LangChain-powered capabilities as robust, machine-readable services.
- **Decoupled Components**: Leveraging LangChain's modularity, the system easily separates Data Indexing from Inference logic.

---

## Key Features

- **Agentic Loop Reasoning**: Multi-turn "thinking" loop that enables complex task execution through tool usage.
- **Conversation Context Memory**: Built-in session-based chat history for multi-turn dialogue persistence.
- **Decoupled SBI Layer**: Professional project structure with clear separation between transmission, logic, and data.
- **High-Performance RAG**: Powered by `langchaingo` with support for multiple document formats (**PDF**, **Markdown**, **Plain Text**).
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
│   └── app/                 # Main entry point (Wiring SBI + Processor + RAG)
├── internal/
│   ├── sbi/                 # Service-Based Interface (API Routing & Handlers)
│   │   ├── processor/       # Internal Business Logic & Agentic reasoning
│   │   └── server.go        # Web Server & Static Asset Mounting
│   ├── service/             # Application Business Logic Layer
│   ├── models/              # Data Structures & Domain Models
│   ├── promptmgr/           # System Prompt Management
│   │   ├── promptmgr.go     # System prompt manager
│   │   └── system.txt       # System prompt (in-memory prompt)
│   ├── tools/               # Dynamic Tool Framework (Add tools here)
│   ├── rag/                 # Low-level RAG logic (Vector DB, Embedding, LLM)
│   ├── history/             # Session-based chat history management
│   └── config/              # Configuration management
├── pkg/
│   └── logger/              # Structured Logging utility (logrus)
├── web/                     # Decoupled Web Frontend (HTML, JS, Assets)
├── documents/               # Local knowledge base (.txt, .md, .pdf)
├── Dockerfile               # Multi-stage production build
└── docker-compose.yml       # Service orchestration
```

---

## Quick Demo (Optional AI Agent features)
This template includes a pre-built demo showcasing **Agentic RAG + Tool Calling** (Office Routing Agent).
- **To Run Demo**: `./demo.sh` (This activates demo files, specialized prompt, and starts Docker).
- **To Clean Demo**: `./demo.sh clean` (This removes demo files and reverts to clean state).

### Testing Steps
1. **Access UI**: Open [http://localhost:8080/ui/](http://localhost:8080/ui/) in your browser.
2. **Indexing**: Click the **"Index"** button (or call `POST /v1/index`) to load `company_rules.md` into the vector database.
3. **Ask a Question**: Type your query in the chat interface.
4. **Clear History**: Click the **"Clear Chat"** button in the sidebar to reset the session history.

### Sample Queries & Expected Result
- **Query**: "I'm planning to take parental leave. Who should I get in touch with? Is he/she free to take a call at the moment?"
- **Expected Result**: 
  - AI identifies **Carol (HR Manager)** from the knowledge base.
  - AI calls the `CheckAvailability` tool to get a random real-time status.
  - AI answers with Carol's contact info and her current availability.

### AI Agent Capabilities (Office Routing Demo)
When the **Quick Demo** is installed, the AI Agent gains the following capabilities:
- **Intelligent Routing (RAG)**: Automatically identifies the correct person to contact based on `company_rules.md` in the knowledge base.
- **Real-time Status Check (Tool)**: Uses a specialized tool to check the "real-time" availability/busy status of employees before answering.
- **Specialized Reasoning**: Follows strict office routing protocols defined in the embedded system prompt.

---

## API Endpoints (SBI)
- `POST /v1/query`: Submit a question to the Processor (supports RAG + Agent tools).
- `DELETE /v1/history`: Clear chat history for a specific session ID.
- `POST /v1/index`: Trigger background re-indexing of the knowledge base.
- `GET /v1/config`: Retrieve dynamic UI branding and system status.
- `GET /v1/tools`: Retrieve the list of currently registered AI Agent tools.

---

## License
This project is licensed under the [MIT License](LICENSE).
