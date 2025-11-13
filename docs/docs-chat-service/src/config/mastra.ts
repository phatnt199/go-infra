import * as dotenv from "dotenv";
import { google } from "@ai-sdk/google";
import { Mastra } from "@mastra/core//mastra";
import { LibSQLVector } from "@mastra/libsql";
import { embed } from "ai";

dotenv.config();

// Set the environment variable that the AI SDK expects
if (process.env.GOOGLE_API_KEY && !process.env.GOOGLE_GENERATIVE_AI_API_KEY) {
	process.env.GOOGLE_GENERATIVE_AI_API_KEY = process.env.GOOGLE_API_KEY;
}

if (!process.env.GOOGLE_GENERATIVE_AI_API_KEY) {
	throw new Error(
		"GOOGLE_GENERATIVE_AI_API_KEY or GOOGLE_API_KEY is required in environment variables"
	);
}

const databaseUrl = process.env.DATABASE_URL || "file:./data/docs-chat.db";

// Initialize LibSQL vector database
export const vectorDb = new LibSQLVector({
	connectionUrl: databaseUrl,
	authToken: process.env.DATABASE_AUTH_TOKEN,
});

// Initialize embedding model
export const embeddingModel = google.textEmbeddingModel("gemini-embedding-001");

// Initialize LLM model for chat
export const chatModel = google("gemini-2.5-flash");

// Initialize Mastra instance
export const mastra = new Mastra({
	vectors: {
		libsql: vectorDb,
	},
});

// Helper function to generate embeddings
export async function generateEmbedding(text: string): Promise<number[]> {
	const { embedding } = await embed({
		model: embeddingModel,
		value: text,
	});
	return embedding;
}

export default mastra;
