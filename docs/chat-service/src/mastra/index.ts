import { Mastra } from "@mastra/core/mastra";
import { LibSQLVector } from "@mastra/libsql";
import * as dotenv from "dotenv";
import { docsAgent } from "./agents/go-infra-docs-agent";

dotenv.config();

const vectorDb = new LibSQLVector({
	connectionUrl: process.env.DATABASE_URL || "file:./data/docs-chat.db",
});

export const mastra = new Mastra({
	agents: {
		docsAgent,
	},
	vectors: {
		libsql: vectorDb,
	},
	server: {
		port: parseInt(process.env.PORT || "3001", 10),
		timeout: 60000, // 60 seconds for RAG queries
		cors: {
			origin: process.env.ALLOWED_ORIGINS?.split(",").map((o) => o.trim()) || [
				"http://localhost:3000",
			],
			allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
			allowHeaders: ["Content-Type", "Authorization"],
			credentials: true,
		},
	},
});
