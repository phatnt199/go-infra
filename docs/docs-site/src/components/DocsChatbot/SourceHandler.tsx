import React, { useState } from "react";
import styles from "./SourceHandler.module.css";

interface Source {
	title: string;
	url: string;
	content: string;
}

interface SourceHandlerProps {
	sources?: Source[];
}

/**
 * Handles source references with smart display and navigation
 * - Only shows sources when they exist
 * - Allows clicking to navigate to the source documentation
 * - Extracts and normalizes documentation links
 */
export default function SourceHandler({ sources }: SourceHandlerProps) {
	const [isExpanded, setIsExpanded] = useState(false);

	if (!sources || sources.length === 0) {
		return null;
	}

	// Filter and deduplicate sources by URL
	const uniqueSources = Array.from(
		new Map(sources.map((s) => [s.url, s])).values()
	);

	const handleSourceClick = (url: string) => {
		// Navigate to the documentation page
		// URLs are already formatted correctly (e.g., "/docs/getting-started")
		window.location.href = url;
	};

	return (
		<div className={styles.sourceHandler}>
			<button
				className={styles.sourceToggle}
				onClick={() => setIsExpanded(!isExpanded)}
				aria-expanded={isExpanded}
			>
				<span className={styles.sourceIcon}>📖</span>
				<span className={styles.sourceLabel}>
					Sources ({uniqueSources.length})
				</span>
				<span
					className={`${styles.arrow} ${isExpanded ? styles.arrowOpen : ""}`}
				>
					▼
				</span>
			</button>

			{isExpanded && (
				<div className={styles.sourcesList}>
					{uniqueSources.map((source, idx) => (
						<div
							key={idx}
							className={styles.sourceItem}
							onClick={() => handleSourceClick(source.url)}
							role="button"
							tabIndex={0}
							onKeyDown={(e) => {
								if (e.key === "Enter" || e.key === " ") {
									handleSourceClick(source.url);
								}
							}}
						>
							<div className={styles.sourceItemHeader}>
								<strong className={styles.sourceTitle}>{source.title}</strong>
								<span className={styles.sourceLink}>{source.url}</span>
								<span className={styles.clickHint}>→</span>
							</div>
							{source.content && (
								<p className={styles.sourceContent}>
									{source.content.substring(0, 150)}...
								</p>
							)}
						</div>
					))}
				</div>
			)}
		</div>
	);
}
