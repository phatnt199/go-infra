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
				"core-concepts/utilities",
			],
		},
		{
			type: "category",
			label: "Core Runtime",
			items: [
				"core-runtime/runtime-overview",
				"core-runtime/events",
				"core-runtime/cqrs",
				"core-runtime/messaging",
				"core-runtime/metadata",
				"core-runtime/serialization",
			],
		},
		{
			type: "category",
			label: "Event Sourcing",
			items: [
				"event-sourcing/overview",
				"event-sourcing/configuration",
				"event-sourcing/projections",
				"event-sourcing/testing",
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
			items: [
				"components/logger",
				"components/crypto",
				"components/health",
				"components/mapper",
				"components/typemapper",
				"components/reflection",
				"components/validator",
				"components/testing",
				"components/otel",
			],
		},
		{
			type: "category",
			label: "Infrastructure",
			items: [
				"infra/postgres",
				"infra/redis",
				"infra/cache",
				"infra/queue",
				"infra/storage",
			],
		},
		{
			type: "category",
			label: "Migration",
			items: ["migration/overview", "migration/goose", "migration/gomigrate"],
		},
		{
			type: "category",
			label: "Usage Patterns",
			items: [
				"usage/direct-vs-fx",
				"usage/customization",
				"usage/module-testing",
				"usage/config-patterns",
			],
		},
		{
			type: "category",
			label: "Project Structures",
			items: [
				"project-structures/minimal",
				"project-structures/modular",
				"project-structures/microservices",
			],
		},
		{
			type: "category",
			label: "Reference",
			items: [
				"reference/configuration-options",
				"reference/fx-modules",
				"reference/public-apis",
			],
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
