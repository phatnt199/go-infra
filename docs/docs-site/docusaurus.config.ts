import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
	title: "go-infra",
	tagline: "Production-ready infrastructure framework for Go applications",
	favicon: "img/favicon.ico",

	// Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
	future: {
		v4: true, // Improve compatibility with the upcoming Docusaurus v4
	},

	// Set the production url of your site here
	url: "https://phatnt199.github.io",
	// Set the /<baseUrl>/ pathname under which your site is served
	// For GitHub pages deployment, it is often '/<projectName>/'
	baseUrl: "/",

	// GitHub pages deployment config.
	// If you aren't using GitHub pages, you don't need these.
	organizationName: "phatnt199", // Usually your GitHub org/user name.
	projectName: "go-infra", // Usually your repo name.

	onBrokenLinks: "throw",

	// Even if you don't use internationalization, you can use this field to set
	// useful metadata like html lang. For example, if your site is Chinese, you
	// may want to replace "en" with "zh-Hans".
	i18n: {
		defaultLocale: "en",
		locales: ["en"],
	},

	presets: [
		[
			"classic",
			{
				docs: {
					sidebarPath: "./sidebars.ts",
					routeBasePath: "docs",
					editUrl: "https://github.com/phatnt199/go-infra/tree/main/docs-site/",
				},
				// Blog is disabled because the site does not use blogging features
				blog: false,
				theme: {
					customCss: "./src/css/custom.css",
				},
			} satisfies Preset.Options,
		],
	],

	themeConfig: {
		// Replace with your project's social card
		image: "img/docusaurus-social-card.jpg",
		colorMode: {
			respectPrefersColorScheme: true,
		},
		navbar: {
			title: "go-infra",
			logo: {
				alt: "go-infra Logo",
				src: "img/logo.svg",
			},
			items: [
				{
					type: "docSidebar",
					sidebarId: "docs",
					position: "left",
					label: "Documentation",
				},
				{
					to: "/docs/examples/users-api",
					position: "left",
					label: "Examples",
				},
				{
					href: "https://github.com/phatnt199/go-infra",
					label: "GitHub",
					position: "right",
				},
			],
		},
		footer: {
			style: "dark",
			links: [
				{
					title: "Documentation",
					items: [
						{
							label: "Getting Started",
							to: "/docs",
						},
						{
							label: "Quick Start",
							to: "/docs/getting-started/quick-start",
						},
						{
							label: "Examples",
							to: "/docs/examples/users-api",
						},
					],
				},
				{
					title: "Community",
					items: [
						{
							label: "GitHub",
							href: "https://github.com/phatnt199/go-infra",
						},
						{
							label: "GitHub Issues",
							href: "https://github.com/phatnt199/go-infra/issues",
						},
						{
							label: "GitHub Discussions",
							href: "https://github.com/phatnt199/go-infra/discussions",
						},
					],
				},
				{
					title: "More",
					items: [
						{
							label: "Go Documentation",
							href: "https://golang.org/doc/",
						},
						{
							label: "Fiber Framework",
							href: "https://gofiber.io",
						},
					],
				},
			],
			copyright: `Copyright © ${new Date().getFullYear()} go-infra. Built with Docusaurus.`,
		},
		prism: {
			theme: prismThemes.github,
			darkTheme: prismThemes.dracula,
			additionalLanguages: ["go", "bash", "json", "yaml"],
		},
	} satisfies Preset.ThemeConfig,
};

export default config;
