import { createTool } from "@mastra/core/tools";
import { z } from "zod";
import { embedMany } from "ai";
import { ModelRouterEmbeddingModel } from "@mastra/core";

// Set the environment variable that the AI SDK expects
if (process.env.GOOGLE_API_KEY && !process.env.GOOGLE_GENERATIVE_AI_API_KEY) {
	process.env.GOOGLE_GENERATIVE_AI_API_KEY = process.env.GOOGLE_API_KEY;
}

const embeddingModel = new ModelRouterEmbeddingModel(
	"google/gemini-embedding-001"
);

export const ragSearchTool = createTool({
	id: "rag-search",
	description: `Search the go-infra documentation to find relevant information. 
  Use this tool when you need to answer questions about:
  - Go-infra features, APIs, and usage
  - Configuration and setup
  - Best practices and examples
  - Troubleshooting and debugging
  
  This tool returns relevant documentation chunks with their source URLs from the documentation site.`,
	inputSchema: z.object({
		query: z
			.string()
			.describe("The search query or question to find relevant documentation"),
		topK: z
			.number()
			.optional()
			.default(5)
			.describe("Number of relevant documents to retrieve (default: 5)"),
	}),
	outputSchema: z.object({
		results: z.array(
			z.object({
				content: z.string().describe("The relevant documentation content"),
				title: z.string().describe("The title of the document"),
				url: z.string().describe("The URL to the documentation page"),
				score: z.number().describe("The relevance score (0-1)"),
			})
		),
		message: z.string().describe("Summary message about the search results"),
	}),
	execute: async ({ context, mastra }) => {
		try {
			const { query, topK } = context;

			// Get vector database from Mastra instance
			const vectorDb = mastra?.getVectors()?.libsql;
			if (!vectorDb) {
				throw new Error("Vector database not initialized in Mastra instance");
			}

			// Generate embedding for the query
			const { embeddings } = await embedMany({
				model: embeddingModel,
				values: [query],
			});

			const queryEmbedding = embeddings[0];

			// Search for relevant documents
			const searchResults = await vectorDb.query({
				indexName: "docs",
				queryVector: queryEmbedding,
				topK: topK * 2, // Get more results to filter
			});

			// Filter by relevance threshold
			const relevanceThreshold = 0.6;
			const relevantResults = searchResults
				.filter((result: any) => result.score >= relevanceThreshold)
				.slice(0, topK);

			if (relevantResults.length === 0) {
				return {
					results: [],
					message: `No relevant documentation found for: "${query}". Try rephrasing your question or ask about specific go-infra features.`,
				};
			}

			// Format results with source URLs
			const results = relevantResults.map((result: any) => {
				const metadata = result.metadata || {};
				const source = metadata.source || "";
				const title = metadata.title || "Unknown Document";
				const content = metadata.content || "";

				// Convert file path to URL
				// Example: "docs/getting-started.mdx" -> "/docs/getting-started"
				let url = "/";
				if (source) {
					// Remove file extension
					const pathWithoutExt = source.replace(/\.(md|mdx)$/, "");
					// Ensure it starts with /
					url = pathWithoutExt.startsWith("/")
						? pathWithoutExt
						: `/${pathWithoutExt}`;
				}

				return {
					content,
					title,
					url,
					score: result.score,
				};
			});

			return {
				results,
				message: `Found ${results.length} relevant documentation sections for: "${query}"`,
			};
		} catch (error: any) {
			console.error("Error in RAG search:", error);
			throw new Error(`Failed to search documentation: ${error.message}`);
		}
	},
});
