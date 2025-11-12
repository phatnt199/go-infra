import { vectorDb, generateEmbedding, chatModel } from "../config/mastra";
import { generateText } from "ai";

interface ChatMessage {
	role: "user" | "assistant" | "system";
	content: string;
}

interface RAGResponse {
	answer: string;
	sources: Array<{
		source: string;
		title: string;
		score: number;
		content: string;
	}>;
}

/**
 * RAG workflow for answering questions about documentation
 */
export async function ragQuery(
	query: string,
	conversationHistory: ChatMessage[] = []
): Promise<RAGResponse> {
	try {
		// Step 1: Generate embedding for the user query
		console.log("Generating query embedding for:", query);
		const queryEmbedding = await generateEmbedding(query);

		// Step 2: Search for relevant documents in vector database
		console.log("Searching for relevant documents...");
		const searchResults = await vectorDb.query({
			indexName: "docs",
			queryVector: queryEmbedding,
			topK: 10, // Retrieve top 10, then filter by relevance
		});

		// Step 3: Filter results by relevance threshold to avoid irrelevant sources
		const relevanceThreshold = 0.6;
		const relevantResults = searchResults.filter(
			(result: any) => result.score >= relevanceThreshold
		);

		console.log(
			`Query: "${query}" | Total results: ${searchResults.length} | Relevant (>${relevanceThreshold}): ${relevantResults.length}`
		);

		// If no relevant results found, inform the user
		if (relevantResults.length === 0) {
			console.warn(
				`No relevant documentation found for query: "${query}" (threshold: ${relevanceThreshold})`
			);
			return {
				answer:
					"I couldn't find relevant documentation to answer your question. Please try rephrasing or ask about specific go-infra features like services, middleware, configuration, or deployment.",
				sources: [],
			};
		}

		// Step 4: Format retrieved context only from relevant results
		const context = relevantResults
			.map((result: any, index: number) => {
				const content = result.metadata?.content || "";
				const title = result.metadata?.title || "Unknown";
				return `[Document ${index + 1}] (Source: ${title}, Score: ${(
					result.score * 100
				).toFixed(1)}%)\n${content}`;
			})
			.join("\n\n---\n\n");

		// Step 5: Build the prompt with context
		const systemPrompt = `You are a helpful AI assistant for the go-infra documentation. 
Your role is to answer questions based on the provided documentation context.

Guidelines:
- Only use information from the provided context to answer questions
- If the context doesn't contain enough information, say so
- Be concise but comprehensive
- Use code examples from the context when relevant
- Format your responses using markdown
- Reference the source documents when possible
- Keep responses complete and well-formed

Context from documentation:
${context}`;

		// Step 6: Generate response using LLM with timeout
		console.log("Generating response from LLM...");
		const timeoutPromise = new Promise<never>(
			(_, reject) =>
				setTimeout(() => reject(new Error("Response timeout")), 30000) // 30 second timeout
		);

		const messages: ChatMessage[] = [
			{ role: "system", content: systemPrompt },
			...conversationHistory,
			{ role: "user", content: query },
		];

		const generatePromise = generateText({
			model: chatModel,
			messages: messages.map((m) => ({
				role: m.role,
				content: m.content,
			})),
			temperature: 0.7,
			maxTokens: 2000, // Increased from 1000 to prevent response cutoff
		});

		const { text: answer } = await Promise.race([
			generatePromise,
			timeoutPromise,
		]);

		// Validate response
		if (!answer || answer.trim().length === 0) {
			throw new Error("Empty response received from model");
		}

		console.log(
			`Response generated successfully (${answer.length} characters)`
		);

		// Step 7: Format sources from relevant results only
		const sources = relevantResults.map((result: any) => ({
			source: result.metadata?.source || "unknown",
			title: result.metadata?.title || "Unknown",
			score: result.score,
			content: (result.metadata?.content || "").substring(0, 200) + "...",
		}));

		return {
			answer,
			sources,
		};
	} catch (error: any) {
		console.error("Error in RAG query:", error);
		if (error.message === "Response timeout") {
			throw new Error(
				"Response generation took too long. Please try again with a simpler question."
			);
		}
		throw new Error("Failed to process query: " + error.message);
	}
}

/**
 * Simple chat without RAG (fallback)
 */
export async function simpleChat(
	query: string,
	conversationHistory: ChatMessage[] = []
): Promise<string> {
	const messages: ChatMessage[] = [
		{
			role: "system",
			content:
				"You are a helpful AI assistant for the go-infra documentation. Answer questions to the best of your ability.",
		},
		...conversationHistory,
		{ role: "user", content: query },
	];

	const { text } = await generateText({
		model: chatModel,
		messages: messages.map((m) => ({
			role: m.role,
			content: m.content,
		})),
		temperature: 0.7,
		maxTokens: 1000,
	});

	return text;
}
