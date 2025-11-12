import * as fs from "fs";
import * as path from "path";
import matter from "gray-matter";
import MarkdownIt from "markdown-it";
import { glob } from "glob";
import { vectorDb, generateEmbedding } from "../config/mastra";

const md = new MarkdownIt();

interface DocumentChunk {
	id: string;
	content: string;
	metadata: {
		source: string;
		title: string;
		section?: string;
		type: "doc" | "blog";
	};
}

// Chunk size configuration (in characters)
const CHUNK_SIZE = 1000;
const CHUNK_OVERLAP = 200;

/**
 * Read and parse markdown/mdx files
 */
function readMarkdownFile(filePath: string): {
	content: string;
	metadata: any;
} {
	const fileContent = fs.readFileSync(filePath, "utf-8");
	const { data, content } = matter(fileContent);
	return { content, metadata: data };
}

/**
 * Chunk text with overlap for better context preservation
 */
function chunkText(text: string, chunkSize: number, overlap: number): string[] {
	const chunks: string[] = [];
	let start = 0;

	while (start < text.length) {
		const end = Math.min(start + chunkSize, text.length);
		const chunk = text.slice(start, end);
		chunks.push(chunk.trim());
		start += chunkSize - overlap;
	}

	return chunks;
}

/**
 * Process a markdown file and create document chunks
 */
function processMarkdownFile(
	filePath: string,
	docsBasePath: string,
	type: "doc" | "blog"
): DocumentChunk[] {
	const { content, metadata } = readMarkdownFile(filePath);

	// Remove markdown syntax for cleaner text
	const plainText = md
		.render(content)
		.replace(/<[^>]*>/g, " ") // Remove HTML tags
		.replace(/\s+/g, " ") // Normalize whitespace
		.trim();

	const chunks = chunkText(plainText, CHUNK_SIZE, CHUNK_OVERLAP);
	const relativePath = path.relative(docsBasePath, filePath);

	return chunks.map((chunk, index) => ({
		id: `${relativePath}#chunk-${index}`,
		content: chunk,
		metadata: {
			source: relativePath,
			title: metadata.title || path.basename(filePath, path.extname(filePath)),
			section: metadata.section,
			type,
		},
	}));
}

/**
 * Embed documents into the vector database
 */
async function embedDocuments(chunks: DocumentChunk[]) {
	console.log(`Embedding ${chunks.length} document chunks...`);

	try {
		// Create index first (3072 matches gemini-embedding-001 model)
		console.log("Creating vector index 'docs'...");
		await vectorDb.createIndex({
			indexName: "docs",
			dimension: 3072,
			metric: "cosine",
		});
		console.log("✅ Index created successfully");
	} catch (error: any) {
		// Index might already exist, which is fine
		if (!error.message?.includes("already exists")) {
			throw error;
		}
		console.log("ℹ️  Index 'docs' already exists, proceeding with upsert...");
	}

	const vectors: number[][] = [];
	const metadata: Record<string, any>[] = [];
	const ids: string[] = [];

	for (let i = 0; i < chunks.length; i++) {
		const chunk = chunks[i];
		console.log(
			`Processing chunk ${i + 1}/${chunks.length}: ${chunk.metadata.source}`
		);

		try {
			const embedding = await generateEmbedding(chunk.content);

			vectors.push(embedding);
			metadata.push({
				...chunk.metadata,
				content: chunk.content,
			});
			ids.push(chunk.id);

			// Add small delay to avoid rate limiting
			await new Promise((resolve) => setTimeout(resolve, 100));
		} catch (error) {
			console.error(`Error embedding chunk ${chunk.id}:`, error);
		}
	}

	// Insert documents into vector database using upsert
	console.log(`Inserting ${vectors.length} documents into vector database...`);
	await vectorDb.upsert({
		indexName: "docs",
		vectors,
		metadata,
		ids,
	});

	console.log("✅ Successfully embedded all documents!");
}

/**
 * Main function to embed all documentation
 */
async function main() {
	const docsBasePath = path.join(__dirname, "../../../go-infra/docs-site");

	console.log("📚 Starting documentation embedding process...");
	console.log(`Base path: ${docsBasePath}`);

	// Find all markdown and mdx files in docs and blog directories
	const docFiles = await glob("docs/**/*.{md,mdx}", { cwd: docsBasePath });
	const blogFiles = await glob("blog/**/*.{md,mdx}", { cwd: docsBasePath });

	console.log(`Found ${docFiles.length} documentation files`);
	console.log(`Found ${blogFiles.length} blog files`);

	// Process all files
	let allChunks: DocumentChunk[] = [];

	for (const file of docFiles) {
		const fullPath = path.join(docsBasePath, file);
		const chunks = processMarkdownFile(fullPath, docsBasePath, "doc");
		allChunks = allChunks.concat(chunks);
	}

	for (const file of blogFiles) {
		const fullPath = path.join(docsBasePath, file);
		const chunks = processMarkdownFile(fullPath, docsBasePath, "blog");
		allChunks = allChunks.concat(chunks);
	}

	console.log(`Total chunks created: ${allChunks.length}`);

	// Embed all chunks
	await embedDocuments(allChunks);
}

// Run the script
main().catch(console.error);
