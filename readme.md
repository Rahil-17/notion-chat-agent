# Chat with Notion (Go)

This Go project demonstrates how to connect a Notion workspace to an LLM (like GPT-3.5) via API to enable natural language queries and summaries from Notion pages. It acts as a lightweight AI agent that reads structured content blocks from Notion and generates context-aware responses.

## Table of Contents

1. [What Is This Project?](#what-is-this-project)
2. [Why Use This Integration?](#why-use-this-integration)
3. [How It Works](#how-it-works)
4. [Trade-Offs and Design Decisions](#trade-offs-and-design-decisions)
5. [Assumptions and Abstractions](#assumptions-and-abstractions)
6. [Running the Code](#running-the-code)
7. [Key Features](#key-features)
8. [Future Improvements](#future-improvements)

---

### What Is This Project?

It’s a simple "Hello World" AI agent that lets you query your Notion pages with natural language by combining:

- Notion API (for content retrieval)
- OpenAI API (for summarization/response generation)
- Go (for speed, simplicity, and portability)

---

### Why Use This Integration?

1. **Turn Notes into Knowledge**:
   - Chat with your second brain (Notion) like a human assistant.

2. **Quick Experimentation with AI Agents**:
   - Useful scaffolding for building more advanced context-aware AI systems.

3. **Portable and Lightweight**:
   - No complex framework — just Go, a few HTTP calls, and smart JSON parsing.

---

### How It Works

#### Core Steps:

1. Fetch page content using Notion API's `blocks/{block_id}/children`.
2. Extract readable text using block type introspection.
3. Feed the text to an LLM via OpenAI API for summarization or answer generation.
4. Print the AI-generated result.

#### Code Highlights:

- `loadEnv()`: Reads required environment variables.
- `fetchNotionBlocks(pageID)`: Pulls all visible block content from a Notion page.
- `extractText(blocks)`: Extracts and flattens `rich_text` content for LLM input.
- `askLLM(question, context)`: Sends a question with Notion page text to GPT-3.5.
- `main()`: Orchestrates the end-to-end flow.

---

### Trade-Offs and Design Decisions

#### Parsing Notion API JSON
Instead of a full SDK, we use `map[string]interface{}` to keep the code minimal and avoid pulling in heavy libraries.

#### GPT-3.5 for Speed and Cost
Used GPT-3.5 via OpenAI API for fast and low-cost responses. For larger or structured pages, token limits can be a bottleneck.

#### No Recursion Yet
Nested blocks like toggles or tables are not parsed recursively. This simplifies the logic but limits deeper extraction.

---

### Assumptions and Abstractions

- The Notion page is **shared with the integration**.
- Content is **mostly paragraph/heading blocks**.
- Environment variables are set for:
  - `NOTION_API_KEY`
  - `NOTION_PAGE_ID`
  - `OPENAI_API_KEY`
- Only supports summarization or Q&A using one full page.

---

### Running the Code

#### 1. Prerequisites

- Go 1.18+ installed
- A Notion integration created and added to the page
- An OpenAI API key with available quota

#### 2. Setup `.env` file

```env
NOTION_API_KEY=secret_xxx
NOTION_PAGE_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
OPENAI_API_KEY=sk-xxx
```

#### 3. Run the Program

```bash
go run main.go
```

It will:
- Fetch your Notion content
- Ask a sample question like "Summarize this page"
- Print the LLM-generated result

---

### Key Features

- Simple integration with Notion + OpenAI
- No third-party SDKs — uses native Go HTTP + JSON
- Flattened rich-text extraction for easy parsing
- Modular design for future extensibility

---

### Future Improvements

1. **Recursive Block Extraction**:
   - Handle nested content like toggles, columns, and tables.

2. **Multiple Pages and Context Windows**:
   - Support summarizing multiple linked pages or databases.

3. **Embeddings-Based Search**:
   - Add vector-based querying for more scalable interactions.

4. **CLI or Web UI**:
   - Add frontend to interact with your Notion via browser or shell.

5. **Replace OpenAI with Free LLMs**:
   - Add support for Groq, Together.ai, or local models for zero-cost operation.