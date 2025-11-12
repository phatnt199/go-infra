# Docs Chat Service

A RAG-powered chatbot service for the go-infra documentation using Mastra, Gemini AI, and LibSQL.

## Features

- 🤖 **AI-Powered**: Uses Google Gemini 2.0 Flash for intelligent responses
- 📚 **RAG System**: Retrieves relevant documentation context before answering
- 🔍 **Vector Search**: Uses Gemini text-embedding-004 for semantic search
- 💾 **Embedded Database**: LibSQL vector database runs alongside the service
- 🎯 **Context-Aware**: Maintains conversation history for better responses
- 📖 **Source Citations**: Shows which documentation pages were used

## Architecture

### Embedding Strategy

**Recursive Character Text Splitter with Overlap**:

- Chunk size: 1000 characters
- Overlap: 200 characters
- Preserves context between chunks
- Ideal for documentation with headings and code blocks

### RAG Workflow

1. **Query Embedding**: User query is embedded using Gemini text-embedding-004
2. **Vector Search**: Top 5 most relevant document chunks are retrieved
3. **Context Building**: Retrieved chunks are formatted with source information
4. **LLM Generation**: Gemini 2.0 Flash generates response with context
5. **Source Attribution**: Returns answer with citations

## Setup

### 1. Install Dependencies

```bash
cd docs-chat-service
npm install
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` and add your Google API key:

```env
GOOGLE_API_KEY=your_actual_google_api_key
PORT=3001
NODE_ENV=development
ALLOWED_ORIGINS=http://localhost:3000
DATABASE_URL=file:./data/docs-chat.db
```

### 3. Embed Documentation

Run the embedding script to process all documentation files:

```bash
npm run embed
```

This will:

- Read all `.md` and `.mdx` files from `../go-infra/docs-site/docs` and `../go-infra/docs-site/blog`
- Split them into chunks with overlap
- Generate embeddings using Gemini
- Store in LibSQL vector database

### 4. Start the Service

```bash
# Development mode (with auto-reload)
npm run dev

# Production mode
npm run build
npm start
```

The service will be available at `http://localhost:3001`.

## API Endpoints

### Health Check

```bash
GET /health
```

Response:

```json
{
	"status": "ok",
	"message": "Docs chat service is running"
}
```

### Chat with RAG

```bash
POST /api/chat
Content-Type: application/json

{
  "message": "How do I use the HTTP server?",
  "history": [
    {
      "role": "user",
      "content": "previous question"
    },
    {
      "role": "assistant",
      "content": "previous answer"
    }
  ]
}
```

Response:

```json
{
	"success": true,
	"answer": "To use the HTTP server...",
	"sources": [
		{
			"source": "docs/http-server/getting-started.md",
			"title": "Getting Started with HTTP Server",
			"content": "..."
		}
	]
}
```

### Simple Chat (No RAG)

```bash
POST /api/chat/simple
Content-Type: application/json

{
  "message": "Hello!",
  "history": []
}
```

## Database

The LibSQL database is stored in `./data/docs-chat.db`. This file is created automatically when you run the embedding script.

### Database Schema

The vector database stores:

- **id**: Unique identifier for each chunk
- **content**: Text content of the chunk
- **embedding**: 768-dimensional vector from Gemini
- **metadata**:
  - source: File path
  - title: Document title
  - section: Section name (if available)
  - type: 'doc' or 'blog'

## Project Structure

```
docs-chat-service/
├── src/
│   ├── config/
│   │   └── mastra.ts           # Mastra configuration
│   ├── services/
│   │   └── rag-service.ts      # RAG workflow implementation
│   ├── scripts/
│   │   └── embed-docs.ts       # Documentation embedding script
│   └── index.ts                # Express server
├── data/                       # Database storage (git-ignored)
├── package.json
├── tsconfig.json
├── .env                        # Environment variables (git-ignored)
└── README.md
```

## Troubleshooting

### Error: GOOGLE_API_KEY is required

Make sure you have created a `.env` file and added your Google API key.

### Error: Cannot find module

Run `npm install` to install all dependencies.

### Database errors

Delete the `./data` directory and re-run `npm run embed` to recreate the database.

### Rate limiting

The embedding script includes a 100ms delay between API calls. If you still hit rate limits, increase the delay in `src/scripts/embed-docs.ts`.

## Development

### Type Checking

```bash
npm run typecheck
```

### Building

```bash
npm run build
```

## License

MIT
