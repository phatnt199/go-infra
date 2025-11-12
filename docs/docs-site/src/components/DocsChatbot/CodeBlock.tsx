import React, { useState } from "react";
import styles from "./CodeBlock.module.css";

interface CodeBlockProps {
	code: string;
	language?: string;
}

export default function CodeBlock({ code, language = "" }: CodeBlockProps) {
	const [copied, setCopied] = useState(false);

	const handleCopy = () => {
		navigator.clipboard.writeText(code).then(() => {
			setCopied(true);
			setTimeout(() => setCopied(false), 2000);
		});
	};

	return (
		<div className={styles.codeBlockWrapper}>
			<div className={styles.codeBlockHeader}>
				{language && <span className={styles.language}>{language}</span>}
				<button
					className={styles.copyButton}
					onClick={handleCopy}
					title={copied ? "Copied!" : "Copy code"}
				>
					{copied ? "✓ Copied" : "Copy"}
				</button>
			</div>
			<pre className={styles.codeBlock}>
				<code className={language ? `language-${language}` : ""}>{code}</code>
			</pre>
		</div>
	);
}
