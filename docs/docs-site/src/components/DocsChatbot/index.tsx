import React, { useState, useRef, useEffect } from "react";
import styles from "./DocsChatbot.module.css";
import MarkdownRenderer from "./MarkdownRenderer";
import SourceHandler from "./SourceHandler";

interface Message {
	role: "user" | "assistant";
	content: string;
	sources?: Array<{
		source: string;
		title: string;
		content: string;
	}>;
}

interface DocsChatbotProps {
	isChatOpen?: boolean;
	setIsChatOpen?: (state: boolean) => void;
}

const API_URL =
	typeof window !== "undefined" && (window as any).REACT_APP_CHAT_API_URL
		? (window as any).REACT_APP_CHAT_API_URL
		: "http://localhost:3001";

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
		setInput("");
		setIsLoading(true);

		const userInput = input; // Capture input before clearing
		let retries = 0;
		const maxRetries = 2;

		try {
			// Prepare conversation history (last 5 messages)
			const history = messages.slice(-5).map((msg) => ({
				role: msg.role,
				content: msg.content,
			}));

			// Retry logic with exponential backoff
			let lastError: Error | null = null;
			for (retries = 0; retries <= maxRetries; retries++) {
				try {
					// Create abort controller with 45 second timeout
					const controller = new AbortController();
					const timeoutId = setTimeout(() => controller.abort(), 45000);

					const response = await fetch(`${API_URL}/api/chat`, {
						method: "POST",
						headers: {
							"Content-Type": "application/json",
						},
						body: JSON.stringify({
							message: userInput,
							history,
						}),
						signal: controller.signal,
					});

					clearTimeout(timeoutId);

					if (!response.ok) {
						throw new Error(`HTTP ${response.status}: ${response.statusText}`);
					}

					const data = await response.json();

					if (data.success) {
						// Validate response structure
						if (!data.answer || data.answer.trim().length === 0) {
							throw new Error("Empty response received from server");
						}

						const assistantMessage: Message = {
							role: "assistant",
							content: data.answer,
							sources: data.sources,
						};
						setMessages((prev) => [...prev, assistantMessage]);
						return; // Success, exit the function
					} else {
						throw new Error(
							data.error || data.message || "Failed to get response"
						);
					}
				} catch (error: any) {
					lastError = error;

					// Don't retry on abort or if max retries reached
					if (
						error.name === "AbortError" ||
						retries >= maxRetries ||
						error.message.includes("Empty response")
					) {
						throw error;
					}

					// Wait before retrying (exponential backoff: 1s, 2s)
					if (retries < maxRetries) {
						const delay = Math.pow(2, retries) * 1000;
						console.log(
							`Retry attempt ${retries + 1}/${maxRetries} after ${delay}ms...`
						);
						await new Promise((resolve) => setTimeout(resolve, delay));
					}
				}
			}

			// If we get here, all retries failed
			throw lastError || new Error("Failed to get response after retries");
		} catch (error: any) {
			console.error("Error sending message:", error);

			let errorMessage: Message;
			if (error.name === "AbortError") {
				errorMessage = {
					role: "assistant",
					content:
						"Sorry, the request took too long to process. The server might be busy. Please try again with a simpler question or wait a moment before retrying.",
				};
			} else if (error.message.includes("Failed to fetch")) {
				errorMessage = {
					role: "assistant",
					content:
						"Unable to connect to the chat service. Please check your internet connection and try again.",
				};
			} else if (error.message.includes("Empty response")) {
				errorMessage = {
					role: "assistant",
					content:
						"The server returned an empty response. Please try rephrasing your question.",
				};
			} else {
				errorMessage = {
					role: "assistant",
					content: `Sorry, I encountered an error: ${error.message}. Please try again.`,
				};
			}

			setMessages((prev) => [...prev, errorMessage]);
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
							{message.role === "assistant" && (
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
