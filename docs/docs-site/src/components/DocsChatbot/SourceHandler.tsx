import React, { useState } from "react";
import styles from "./SourceHandler.module.css";
import { useHistory } from "@docusaurus/router";

interface Source {
	source: string;
	title: string;
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

	// Filter and deduplicate sources
	const uniqueSources = Array.from(
		new Map(sources.map((s) => [s.source, s])).values()
	);

	const handleSourceClick = (sourceUrl: string) => {
		// Normalize the URL - if it's a relative path or documentation path
		let navigateUrl = sourceUrl;

		// If it's already a full URL, open in new tab
		if (sourceUrl.startsWith("http://") || sourceUrl.startsWith("https://")) {
			window.open(sourceUrl, "_blank");
			return;
		}

		// If it's a path like "docs/..." or "/docs/...", navigate within site
		if (!navigateUrl.startsWith("/")) {
			navigateUrl = "/" + navigateUrl;
		}

		// Use window.location for in-site navigation or redirect
		if (navigateUrl.includes("docs/") || navigateUrl.includes("/docs/")) {
			window.location.href = navigateUrl;
		} else {
			// Try to navigate using router or fallback to location
			window.location.href = navigateUrl;
		}
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
							onClick={() => handleSourceClick(source.source)}
							role="button"
							tabIndex={0}
							onKeyDown={(e) => {
								if (e.key === "Enter" || e.key === " ") {
									handleSourceClick(source.source);
								}
							}}
						>
							<div className={styles.sourceItemHeader}>
								<strong className={styles.sourceTitle}>{source.title}</strong>
								<span className={styles.sourceLink}>
									{truncateUrl(source.source)}
								</span>
								<span className={styles.clickHint}>→</span>
							</div>
							{source.content && (
								<p className={styles.sourceContent}>{source.content}</p>
							)}
						</div>
					))}
				</div>
			)}
		</div>
	);
}

/**
 * Truncate long URLs for display
 */
function truncateUrl(url: string, maxLength: number = 40): string {
	if (url.length <= maxLength) return url;

	// Try to show the relevant part
	const parts = url.split("/");
	const filename = parts[parts.length - 1];

	if (filename.length > maxLength) {
		return "..." + filename.slice(-maxLength);
	}

	return "..." + url.slice(-maxLength);
}
