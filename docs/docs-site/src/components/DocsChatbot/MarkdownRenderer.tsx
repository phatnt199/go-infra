import React from "react";
import styles from "./MarkdownRenderer.module.css";
import CodeBlock from "./CodeBlock";

interface MarkdownRendererProps {
	content: string;
	isUser?: boolean;
}

/**
 * Parses and renders markdown content with support for:
 * - Code blocks with syntax highlighting (using <pre> and <code>)
 * - Inline code
 * - Bold and italic text
 * - Links
 * - Lists
 * - Headers
 */
export default function MarkdownRenderer({
	content,
	isUser = false,
}: MarkdownRendererProps): React.ReactNode {
	// If user message, just return plain text
	if (isUser) {
		return <>{content}</>;
	}

	// Split content by code blocks first
	const parts = content.split(/(```[\s\S]*?```)/);

	return (
		<div className={styles.markdownContent}>
			{parts.map((part, index) => {
				// Check if this is a code block
				if (part.startsWith("```")) {
					const codeContent = part
						.replace(/^```[\w]*\n?/, "") // Remove opening ```
						.replace(/\n?```$/, ""); // Remove closing ```
					const language = part.match(/^```(\w+)/)?.[1] || "";

					return (
						<CodeBlock key={index} code={codeContent} language={language} />
					);
				}

				// Parse inline markdown
				return (
					<React.Fragment key={index}>
						{parseInlineMarkdown(part)}
					</React.Fragment>
				);
			})}
		</div>
	);
}

/**
 * Parse inline markdown elements like bold, italic, inline code, and links
 */
function parseInlineMarkdown(text: string): React.ReactNode {
	if (!text) return null;

	let nodes: React.ReactNode[] = [text];

	// Process bold text
	nodes = nodes.flatMap((node) => {
		if (typeof node !== "string") return node;
		return parseMarkdownPattern(node, /\*\*(.+?)\*\*/g, (content) => (
			<strong key={`bold-${content}`}>{content}</strong>
		));
	});

	// Process italic text
	nodes = nodes.flatMap((node) => {
		if (typeof node !== "string") return node;
		return parseMarkdownPattern(node, /\*(.+?)\*/g, (content) => (
			<em key={`italic-${content}`}>{content}</em>
		));
	});

	// Process inline code
	nodes = nodes.flatMap((node) => {
		if (typeof node !== "string") return node;
		return parseMarkdownPattern(node, /`([^`]+)`/g, (content) => (
			<code key={`code-${content}`} className={styles.inlineCode}>
				{content}
			</code>
		));
	});

	// Process links
	nodes = nodes.flatMap((node) => {
		if (typeof node !== "string") return node;
		const result: React.ReactNode[] = [];
		let lastIndex = 0;
		let match;
		const regex = /\[([^\]]+)\]\(([^)]+)\)/g;

		while ((match = regex.exec(node)) !== null) {
			if (match.index > lastIndex) {
				result.push(node.substring(lastIndex, match.index));
			}
			result.push(
				<a
					key={`link-${match[1]}-${match[2]}`}
					href={match[2]}
					target="_blank"
					rel="noopener noreferrer"
				>
					{match[1]}
				</a>
			);
			lastIndex = regex.lastIndex;
		}

		if (lastIndex < node.length) {
			result.push(node.substring(lastIndex));
		}

		return result.length > 0 ? result : node;
	});

	return nodes;
}

/**
 * Helper function to parse markdown pattern
 */
function parseMarkdownPattern(
	text: string,
	regex: RegExp,
	render: (content: string) => React.ReactNode
): React.ReactNode[] {
	const result: React.ReactNode[] = [];
	let lastIndex = 0;
	let match;

	regex.lastIndex = 0; // Reset regex state

	while ((match = regex.exec(text)) !== null) {
		if (match.index > lastIndex) {
			result.push(text.substring(lastIndex, match.index));
		}
		result.push(render(match[1]));
		lastIndex = regex.lastIndex;
	}

	if (lastIndex < text.length) {
		result.push(text.substring(lastIndex));
	}

	return result.length > 0 ? result : [text];
}
