import * as fs from "fs";
import * as path from "path";
import { fileURLToPath } from "url";
import matter from "gray-matter";
import { glob } from "glob";
import { MDocument } from "@mastra/rag";
import { embedMany } from "ai";
import { LibSQLVector } from "@mastra/libsql";
import * as dotenv from "dotenv";
import { ModelRouterEmbeddingModel } from "@mastra/core";

dotenv.config();

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Validate required environment variables
if (process.env.GOOGLE_API_KEY && !process.env.GOOGLE_GENERATIVE_AI_API_KEY) {
	process.env.GOOGLE_GENERATIVE_AI_API_KEY = process.env.GOOGLE_API_KEY;
}

if (!process.env.GOOGLE_GENERATIVE_AI_API_KEY) {
	throw new Error(
		"GOOGLE_GENERATIVE_AI_API_KEY or GOOGLE_API_KEY is required in environment variables"
	);
}

// Using 768 dimensions for optimal balance between performance and accuracy
// Gemini embedding model supports flexible dimensions (128-3072)
const EMBEDDING_DIMENSION = 768;
const embeddingModel = new ModelRouterEmbeddingModel(
	"google/gemini-embedding-001"
);

// Initialize LibSQL vector database
let vectorDb: LibSQLVector;

try {
	const dbUrl = process.env.DATABASE_URL || "file:./data/docs-chat.db";

	// If using a file URL (file:./data/...), ensure the directory exists
	let filePath = dbUrl;
	if (dbUrl.startsWith("file:")) {
		filePath = dbUrl.replace(/^file:/, "");
	}

	// If this looks like a local filesystem path, create the directory
	// (skip if it's a remote URL like postgres:// or sqlite+...://)
	if (!/^[a-zA-Z]+:\/\//.test(filePath) && filePath !== "/") {
		const dir = path.dirname(filePath);
		if (dir && !fs.existsSync(dir)) {
			fs.mkdirSync(dir, { recursive: true });
			console.log(`Created directory for DB at ${dir}`);
		}
	}

	vectorDb = new LibSQLVector({
		connectionUrl: dbUrl,
	});
} catch (error) {
	console.error(
		"Failed to initialize LibSQLVector or create DB directory:",
		error
	);
	throw error;
}

interface DocumentChunk {
	id: string;
	text: string;
	metadata: {
		source: string;
		title: string;
		section?: string;
		type: "doc" | "blog";
		url: string; // Add URL for reference
	};
}

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
 * Convert file path to documentation URL
 */
function filePathToUrl(relativePath: string, type: "doc" | "blog"): string {
	// Remove file extension
	const pathWithoutExt = relativePath.replace(/\.(md|mdx)$/, "");

	// For docs: docs/getting-started.mdx -> /docs/getting-started
	// For blog: blog/2024-01-01-post.mdx -> /blog/2024-01-01-post
	const url = pathWithoutExt.startsWith("/")
		? pathWithoutExt
		: `/${pathWithoutExt}`;

	return url;
}

/**
 * Process a markdown file and create document chunks using Mastra RAG
 */
async function processMarkdownFile(
	filePath: string,
	docsBasePath: string,
	type: "doc" | "blog"
): Promise<DocumentChunk[]> {
	const { content, metadata } = readMarkdownFile(filePath);
	const relativePath = path.relative(docsBasePath, filePath);
	const url = filePathToUrl(relativePath, type);

	// Create MDocument from markdown content
	const doc = MDocument.fromMarkdown(content);

	// Optimized chunking: larger chunks with less overlap to reduce total chunk count
	// This helps stay within RPD (requests per day) limits
	const chunks = await doc.chunk({
		strategy: "recursive",
		maxSize: 1500, // Increased from 800 to reduce chunk count
		overlap: 150, // Proportional overlap
	});

	// Map chunks to our format with URL metadata
	return chunks.map((chunk, index) => ({
		id: `${relativePath}#chunk-${index}`,
		text: chunk.text,
		metadata: {
			source: relativePath,
			title: metadata.title || path.basename(filePath, path.extname(filePath)),
			section: metadata.section,
			type,
			url, // Include the documentation URL
		},
	}));
}

/**
 * Load progress from cache file
 */
function loadProgress(): Set<string> {
	const progressFile = path.join(__dirname, "../../data/embed-progress.json");
	try {
		if (fs.existsSync(progressFile)) {
			const data = JSON.parse(fs.readFileSync(progressFile, "utf-8"));
			return new Set(data.processedIds || []);
		}
	} catch (error) {
		console.warn("Could not load progress file:", error);
	}
	return new Set();
}

/**
 * Save progress to cache file
 */
function saveProgress(processedIds: Set<string>) {
	const progressFile = path.join(__dirname, "../../data/embed-progress.json");
	const dir = path.dirname(progressFile);
	if (!fs.existsSync(dir)) {
		fs.mkdirSync(dir, { recursive: true });
	}
	fs.writeFileSync(
		progressFile,
		JSON.stringify(
			{
				processedIds: Array.from(processedIds),
				timestamp: new Date().toISOString(),
			},
			null,
			2
		)
	);
}

/**
 * Embed documents into the vector database with optimized rate limiting
 */
async function embedDocuments(chunks: DocumentChunk[]) {
	console.log(`Embedding ${chunks.length} document chunks...`);

	// Load progress to resume from where we left off
	const processedIds = loadProgress();
	const remainingChunks = chunks.filter((chunk) => !processedIds.has(chunk.id));

	if (remainingChunks.length === 0) {
		console.log("✅ All chunks already processed!");
		return;
	}

	console.log(
		`Found ${processedIds.size} already processed. Processing ${remainingChunks.length} remaining chunks...`
	);

	try {
		// Create index first (using optimized 768 dimensions)
		console.log("Creating vector index 'docs'...");
		await vectorDb.createIndex({
			indexName: "docs",
			dimension: EMBEDDING_DIMENSION,
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

	// Optimized batching strategy for free tier:
	// - Smaller batch size (5 instead of 10) to stay under RPM limit
	// - Longer delays between batches to respect TPM limits
	// - Process max 100 chunks per run to stay under RPD limit
	const batchSize = 5;
	const delayBetweenBatches = 3000; // 3 seconds between batches
	const maxChunksPerRun = 100; // Stay well under RPD limit of 1000

	const chunksToProcess = remainingChunks.slice(0, maxChunksPerRun);
	console.log(
		`Processing ${chunksToProcess.length} chunks (max ${maxChunksPerRun} per run)...`
	);

	for (let i = 0; i < chunksToProcess.length; i += batchSize) {
		const batch = chunksToProcess.slice(i, i + batchSize);
		const batchNum = Math.floor(i / batchSize) + 1;
		const totalBatches = Math.ceil(chunksToProcess.length / batchSize);

		console.log(
			`Processing batch ${batchNum}/${totalBatches} (${batch.length} chunks)...`
		);

		try {
			// Generate embeddings for the batch
			const startTime = Date.now();
			const { embeddings } = await embedMany({
				model: embeddingModel,
				values: batch.map((chunk) => chunk.text),
			});
			const duration = Date.now() - startTime;

			// Collect data for upsert
			embeddings.forEach((embedding, idx) => {
				const chunk = batch[idx];
				vectors.push(embedding);
				metadata.push({
					...chunk.metadata,
					content: chunk.text, // Store the actual text content
				});
				ids.push(chunk.id);
				processedIds.add(chunk.id);
			});

			console.log(`  ✓ Batch completed in ${duration}ms`);

			// Save progress after each successful batch
			saveProgress(processedIds);

			// Rate limiting: wait between batches to respect API limits
			// Free tier: 100 RPM, 30K TPM, 1000 RPD
			if (i + batchSize < chunksToProcess.length) {
				console.log(
					`  ⏳ Waiting ${delayBetweenBatches}ms to respect rate limits...`
				);
				await new Promise((resolve) =>
					setTimeout(resolve, delayBetweenBatches)
				);
			}
		} catch (error: any) {
			console.error(
				`❌ Error embedding batch ${batchNum}:`,
				error.message || error
			);

			// If rate limited, wait longer and continue
			if (error.message?.includes("429") || error.message?.includes("quota")) {
				console.log(
					"⚠️  Rate limit hit. Waiting 60 seconds before continuing..."
				);
				await new Promise((resolve) => setTimeout(resolve, 60000));
			} else {
				// For other errors, save progress and throw
				saveProgress(processedIds);
				throw error;
			}
		}
	} // Insert documents into vector database using upsert
	if (vectors.length > 0) {
		console.log(
			`Inserting ${vectors.length} documents into vector database...`
		);
		await vectorDb.upsert({
			indexName: "docs",
			vectors,
			metadata,
			ids,
		});
		console.log("✅ Successfully embedded and inserted documents!");
	}

	const remaining = remainingChunks.length - chunksToProcess.length;
	if (remaining > 0) {
		console.log(
			`\n⚠️  ${remaining} chunks remaining. Run the script again to continue.`
		);
		console.log(
			`Progress: ${processedIds.size}/${chunks.length} chunks completed`
		);
	} else {
		console.log("\n🎉 All documents successfully embedded!");
		// Clean up progress file when done
		const progressFile = path.join(__dirname, "../../data/embed-progress.json");
		if (fs.existsSync(progressFile)) {
			fs.unlinkSync(progressFile);
			console.log("✓ Progress file cleaned up");
		}
	}
}

/**
 * Main function to embed all documentation
 */
async function main() {
	const docsBasePath = path.join(__dirname, "../../../docs-site");

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
		console.log(`Processing doc: ${file}...`);
		const chunks = await processMarkdownFile(fullPath, docsBasePath, "doc");
		allChunks = allChunks.concat(chunks);
	}

	for (const file of blogFiles) {
		const fullPath = path.join(docsBasePath, file);
		console.log(`Processing blog: ${file}...`);
		const chunks = await processMarkdownFile(fullPath, docsBasePath, "blog");
		allChunks = allChunks.concat(chunks);
	}

	console.log(`Total chunks created: ${allChunks.length}`);

	// Embed all chunks
	await embedDocuments(allChunks);
}

// Run the script
main().catch(console.error);
