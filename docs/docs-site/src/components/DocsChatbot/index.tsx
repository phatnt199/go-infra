import React, { useState, useRef, useEffect } from "react";
import "highlight.js/styles/atom-one-dark.css";
import styles from "./DocsChatbot.module.css";
import MarkdownRenderer from "./MarkdownRenderer";
import SourceHandler from "./SourceHandler";
import { mastraClient } from "../../lib/mastra-client";

interface Message {
	role: "user" | "assistant";
	content: string;
	sources?: Array<{
		title: string;
		url: string;
		content: string;
	}>;
}

interface DocsChatbotProps {
	isChatOpen?: boolean;
	setIsChatOpen?: (state: boolean) => void;
}

export default function DocsChatbot({
	isChatOpen = false,
	setIsChatOpen,
}: DocsChatbotProps) {
	const [isOpen, setIsOpen] = useState(isChatOpen);
	const [messages, setMessages] = useState<Message[]>([
		{
			role: "assistant",
			content:
				"Hi! I'm here to help you with the go-infra documentation. Ask me anything!",
		},
	]);
	const [input, setInput] = useState("");
	const [isLoading, setIsLoading] = useState(false);
	const messagesEndRef = useRef<HTMLDivElement>(null);

	const scrollToBottom = () => {
		messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
	};

	useEffect(() => {
		scrollToBottom();
	}, [messages]);

	useEffect(() => {
		// Sync with parent state
		if (setIsChatOpen) {
			setIsChatOpen(isOpen);
		}
	}, [isOpen, setIsChatOpen]);

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		if (!input.trim() || isLoading) return;

		const userMessage: Message = {
			role: "user",
			content: input,
		};

		setMessages((prev) => [...prev, userMessage]);
		const userInput = input;
		setInput("");
		setIsLoading(true);

		// Create a placeholder message for streaming
		const assistantMessageIndex = messages.length + 1;
		setMessages((prev) => [
			...prev,
			{
				role: "assistant",
				content: "",
			},
		]);

		try {
			// Get the docsAgent from Mastra client
			const agent = mastraClient.getAgent("docsAgent");

			// Prepare conversation history (last 5 messages for context)
			const conversationHistory = messages.slice(-5).map((msg) => ({
				role: msg.role,
				content: msg.content,
			}));

			// Stream the response
			const stream = await agent.stream({
				messages: [
					...conversationHistory,
					{ role: "user", content: userInput },
				],
			} as any);

			let fullText = "";
			let sources: Array<{ title: string; url: string; content: string }> = [];

			// Process the stream
			await stream.processDataStream({
				onChunk: async (chunk: any) => {
					if (chunk.type === "text-delta") {
						fullText += chunk.payload.text;
						// Update the message in real-time
						setMessages((prev) => {
							const newMessages = [...prev];
							newMessages[assistantMessageIndex] = {
								role: "assistant",
								content: fullText,
								sources,
							};
							return newMessages;
						});
					} else if (chunk.type === "tool-result") {
						// Extract sources from RAG tool results
						const toolResult = chunk.payload.result;
						if (toolResult?.results && Array.isArray(toolResult.results)) {
							sources = toolResult.results.map((r: any) => ({
								title: r.title || "Unknown",
								url: r.url || "#",
								content: r.content || "",
							}));
							// Update message with sources
							setMessages((prev) => {
								const newMessages = [...prev];
								newMessages[assistantMessageIndex] = {
									role: "assistant",
									content: fullText,
									sources,
								};
								return newMessages;
							});
						}
					}
				},
			});

			// Ensure final message has all content
			setMessages((prev) => {
				const newMessages = [...prev];
				newMessages[assistantMessageIndex] = {
					role: "assistant",
					content:
						fullText ||
						"I apologize, but I couldn't generate a response. Please try again.",
					sources,
				};
				return newMessages;
			});
		} catch (error: any) {
			console.error("Error sending message:", error);

			// Update the assistant message with error
			setMessages((prev) => {
				const newMessages = [...prev];
				newMessages[assistantMessageIndex] = {
					role: "assistant",
					content: `Sorry, I encountered an error: ${
						error.message || "Unknown error"
					}. Please try again.`,
				};
				return newMessages;
			});
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div className={`${isOpen ? styles.chatExpanded : styles.chatCollapsed}`}>
			{/* Toggle Button - visible in both states */}
			<button
				className={styles.chatToggleButton}
				onClick={() => setIsOpen(!isOpen)}
				aria-label="Toggle documentation assistant"
				title={isOpen ? "Close chat" : "Open chat"}
			>
				<span className={styles.toggleArrow}>{isOpen ? "✕" : "💬"}</span>
			</button>

			{/* Chat Content - hidden when collapsed */}
			<div className={styles.chatContent}>
				<div className={styles.chatHeader}>
					<h3>📚 Documentation Assistant</h3>
				</div>

				<div className={styles.chatMessages}>
					{messages.map((message, index) => (
						<div
							key={index}
							className={`${styles.message} ${
								message.role === "user"
									? styles.userMessage
									: styles.assistantMessage
							}`}
						>
							<div className={styles.messageContent}>
								{message.role === "user" ? (
									message.content
								) : (
									<MarkdownRenderer content={message.content} />
								)}
							</div>
							{message.role === "assistant" && message.sources && (
								<SourceHandler sources={message.sources} />
							)}
						</div>
					))}
					{isLoading && (
						<div className={`${styles.message} ${styles.assistantMessage}`}>
							<div className={styles.messageContent}>
								<div className={styles.loadingDots}>
									<span></span>
									<span></span>
									<span></span>
								</div>
							</div>
						</div>
					)}
					<div ref={messagesEndRef} />
				</div>

				<form className={styles.chatInput} onSubmit={handleSubmit}>
					<input
						type="text"
						value={input}
						onChange={(e) => setInput(e.target.value)}
						placeholder="Ask about the documentation..."
						disabled={isLoading}
						aria-label="Chat message input"
					/>
					<button
						type="submit"
						disabled={isLoading || !input.trim()}
						aria-label="Send message"
					>
						Send
					</button>
				</form>
			</div>
		</div>
	);
}
