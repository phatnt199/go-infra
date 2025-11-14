import React from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeHighlight from "rehype-highlight";
import type { Components } from "react-markdown";
import type { Plugin } from "unified";
import { visit } from "unist-util-visit";
import styles from "./MarkdownRenderer.module.css";
import CodeBlock from "./CodeBlock";

interface MarkdownRendererProps {
	content: string;
	isUser?: boolean;
}

/**
 * Rehype plugin to skip highlighting inline code
 * Only highlights code in <pre><code> (block code), not standalone <code> tags
 */
const skipInlineCodeHighlight: Plugin = () => {
	return (tree: any) => {
		visit(tree, "element", (node: any) => {
			// Skip code highlighting for inline code elements that are NOT inside pre
			if (node.tagName === "code" && node.parent?.tagName !== "pre") {
				// Remove hljs class if it exists
				if (node.properties.className) {
					node.properties.className = node.properties.className.filter(
						(cls: string) =>
							!cls.startsWith("hljs") && !cls.startsWith("language-")
					);
				}
			}
		});
	};
};

/**
 * Custom component overrides for react-markdown
 * Provides syntax highlighting for code blocks and custom styling
 */
const customComponents: Components = {
	code: ({ node, inline, className, children, ...props }: any) => {
		const codeContent = String(children);
		const hasNewline = codeContent.includes("\n");
		const shouldRenderInline = inline ?? !hasNewline;

		if (shouldRenderInline) {
			return (
				<code className={styles.inlineCode} {...props}>
					{codeContent}
				</code>
			);
		}

		// Block code (fenced code blocks with ```)
		const match = /language-(\w+)/.exec(className || "");
		const language = match ? match[1] : "";

		return (
			<CodeBlock code={codeContent.replace(/\n$/, "")} language={language} />
		);
	},
	pre: ({ children, ...props }: any) => {
		// Extract language from the code element if present
		return <>{children}</>;
	},
	a: ({ href, children, ...props }) => (
		<a href={href} target="_blank" rel="noopener noreferrer" {...props}>
			{children}
		</a>
	),
	h1: ({ children, ...props }) => (
		<h1 className={styles.heading1} {...props}>
			{children}
		</h1>
	),
	h2: ({ children, ...props }) => (
		<h2 className={styles.heading2} {...props}>
			{children}
		</h2>
	),
	h3: ({ children, ...props }) => (
		<h3 className={styles.heading3} {...props}>
			{children}
		</h3>
	),
	ul: ({ children, ...props }) => (
		<ul className={styles.list} {...props}>
			{children}
		</ul>
	),
	ol: ({ children, ...props }) => (
		<ol className={styles.orderedList} {...props}>
			{children}
		</ol>
	),
	li: ({ children, ...props }) => (
		<li className={styles.listItem} {...props}>
			{children}
		</li>
	),
	blockquote: ({ children, ...props }) => (
		<blockquote className={styles.blockquote} {...props}>
			{children}
		</blockquote>
	),
	table: ({ children, ...props }) => (
		<table className={styles.table} {...props}>
			{children}
		</table>
	),
	thead: ({ children, ...props }) => (
		<thead className={styles.tableHead} {...props}>
			{children}
		</thead>
	),
	tbody: ({ children, ...props }) => (
		<tbody className={styles.tableBody} {...props}>
			{children}
		</tbody>
	),
	tr: ({ children, ...props }) => (
		<tr className={styles.tableRow} {...props}>
			{children}
		</tr>
	),
	th: ({ children, ...props }) => (
		<th className={styles.tableHeader} {...props}>
			{children}
		</th>
	),
	td: ({ children, ...props }) => (
		<td className={styles.tableCell} {...props}>
			{children}
		</td>
	),
};

/**
 * Renders markdown content with support for:
 * - GitHub Flavored Markdown (GFM) with tables, strikethrough, etc.
 * - Syntax highlighting via rehype-highlight
 * - Code blocks with language support
 * - Inline code, bold, italic, links
 * - Lists, blockquotes, and headers
 */
export default function MarkdownRenderer({
	content,
	isUser = false,
}: MarkdownRendererProps): React.ReactNode {
	// If user message, just return plain text
	if (isUser) {
		return <>{content}</>;
	}

	return (
		<div className={styles.markdownContent}>
			<ReactMarkdown
				remarkPlugins={[remarkGfm]}
				rehypePlugins={[skipInlineCodeHighlight, rehypeHighlight]}
				components={customComponents}
			>
				{content}
			</ReactMarkdown>
		</div>
	);
}
