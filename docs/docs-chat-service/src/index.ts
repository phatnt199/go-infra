import express, { Request, Response } from "express";
import cors from "cors";
import * as dotenv from "dotenv";
import { ragQuery, simpleChat } from "./services/rag-service";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3001;
const allowedOrigins = (process.env.ALLOWED_ORIGINS || "http://localhost:3000")
	.split(",")
	.map((origin) => origin.trim());

console.log("✅ CORS allowed origins:", allowedOrigins);

// Middleware
app.use(
	cors({
		origin: allowedOrigins,
		credentials: true,
		methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
		allowedHeaders: ["Content-Type", "Authorization"],
	})
);
app.use(express.json());

// Health check endpoint
app.get("/health", (req: Request, res: Response) => {
	res.json({ status: "ok", message: "Docs chat service is running" });
});

// Chat endpoint with RAG
app.post("/api/chat", async (req: Request, res: Response) => {
	// Set a reasonable timeout for this request (50 seconds)
	req.setTimeout(50000);

	try {
		const { message, history = [] } = req.body;

		if (!message) {
			return res.status(400).json({ error: "Message is required" });
		}

		if (typeof message !== "string" || message.trim().length === 0) {
			return res
				.status(400)
				.json({ error: "Message must be a non-empty string" });
		}

		console.log(
			`[${new Date().toISOString()}] Received chat request: "${message.substring(
				0,
				50
			)}..."`
		);

		// Process query with RAG with timeout handling
		let response: any;
		try {
			response = await Promise.race([
				ragQuery(message, history),
				new Promise(
					(_, reject) =>
						setTimeout(
							() => reject(new Error("RAG query processing timeout")),
							45000
						) // 45 second timeout for processing
				),
			]);
		} catch (error: any) {
			if (error.message === "RAG query processing timeout") {
				return res.status(504).json({
					success: false,
					error: "Processing timeout",
					message:
						"The response took too long to generate. Please try again with a simpler question.",
				});
			}
			throw error;
		}

		// Validate response before sending
		if (!response || !response.answer || response.answer.trim().length === 0) {
			console.warn("Empty or invalid response generated");
			return res.status(500).json({
				success: false,
				error: "Invalid response",
				message: "The service generated an invalid response. Please try again.",
			});
		}

		console.log(
			`✅ Chat response successful (${response.answer.length} characters)`
		);
		res.json({
			success: true,
			answer: response.answer,
			sources: response.sources || [],
		});
	} catch (error: any) {
		console.error(
			`❌ Error processing chat [${new Date().toISOString()}]:`,
			error
		);
		res.status(500).json({
			success: false,
			error: "Failed to process chat request",
			message: error.message || "An unexpected error occurred",
		});
	}
});

// Simple chat endpoint (without RAG)
app.post("/api/chat/simple", async (req: Request, res: Response) => {
	try {
		const { message, history = [] } = req.body;

		if (!message) {
			return res.status(400).json({ error: "Message is required" });
		}

		const response = await simpleChat(message, history);

		res.json({
			success: true,
			answer: response,
		});
	} catch (error: any) {
		console.error("Error processing simple chat:", error);
		res.status(500).json({
			success: false,
			error: "Failed to process chat request",
			message: error.message,
		});
	}
});

// Start server
app.listen(PORT, () => {
	console.log(`🚀 Docs chat service is running on http://localhost:${PORT}`);
	console.log(`📚 Ready to answer questions about go-infra documentation`);
});

export default app;
