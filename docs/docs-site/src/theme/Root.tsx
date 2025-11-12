import React, { useState } from "react";
import Root from "@theme-original/Root";
import DocsChatbot from "../components/DocsChatbot";
import styles from "./Root.module.css";

export default function RootWrapper(props) {
	const [isChatOpen, setIsChatOpen] = useState(true);

	return (
		<div
			className={`${styles.rootWrapper} ${
				isChatOpen ? styles.rootWrapperChatOpen : ""
			}`}
		>
			<div className={styles.mainContent}>
				<Root {...props} />
			</div>
			<DocsChatbot isChatOpen={isChatOpen} setIsChatOpen={setIsChatOpen} />
		</div>
	);
}
