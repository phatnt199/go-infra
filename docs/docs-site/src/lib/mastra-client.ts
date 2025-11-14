import { MastraClient } from "@mastra/client-js";

// Get the API URL from environment or default to localhost
const API_URL =
	typeof window !== "undefined" && (window as any).REACT_APP_CHAT_API_URL
		? (window as any).REACT_APP_CHAT_API_URL
		: "http://localhost:3001";

export const mastraClient = new MastraClient({
	baseUrl: API_URL,
});
