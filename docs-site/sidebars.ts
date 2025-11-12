import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
	docs: [
		"intro",
		{
			type: "category",
			label: "Getting Started",
			collapsed: false,
			items: ["getting-started/installation", "getting-started/quick-start"],
		},
		{
			type: "category",
			label: "Core Concepts",
			items: [
				"core-concepts/architecture",
				"core-concepts/modules",
				"core-concepts/dependency-injection",
				"core-concepts/configuration",
			],
		},
		{
			type: "category",
			label: "HTTP Server",
			items: ["http-server/building-apis", "http-server/crud-operations"],
		},
		{
			type: "category",
			label: "Database",
			items: ["database/setup", "database/migrations"],
		},
		{
			type: "category",
			label: "Authentication & Security",
			items: ["authentication/getting-started"],
		},
		{
			type: "category",
			label: "Components",
			items: ["components/logger", "components/crypto"],
		},
		{
			type: "category",
			label: "Advanced Topics",
			items: ["advanced/overview"],
		},
		{
			type: "category",
			label: "Deployment",
			items: ["deployment/production"],
		},
		{
			type: "category",
			label: "Examples",
			items: [
				"examples/users-api",
				"examples/authentication-service",
				"examples/microservices",
			],
		},
	],
};

export default sidebars;
