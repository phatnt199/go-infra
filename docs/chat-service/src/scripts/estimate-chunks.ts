import * as fs from "fs";
import * as path from "path";
import { fileURLToPath } from "url";
import matter from "gray-matter";
import { glob } from "glob";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

/**
 * Estimate the number of chunks that will be created
 * This helps understand how many API calls will be needed
 */
async function estimateChunks() {
	const docsBasePath = path.join(__dirname, "../../../docs-site");

	console.log("📊 Estimating chunk count for embedding...\n");
	console.log(`Base path: ${docsBasePath}\n`);

	// Find all markdown and mdx files
	const docFiles = await glob("docs/**/*.{md,mdx}", { cwd: docsBasePath });
	const blogFiles = await glob("blog/**/*.{md,mdx}", { cwd: docsBasePath });

	console.log(`Found ${docFiles.length} documentation files`);
	console.log(`Found ${blogFiles.length} blog files`);
	console.log(`Total files: ${docFiles.length + blogFiles.length}\n`);

	let totalWords = 0;
	let totalChars = 0;
	let fileStats: Array<{
		path: string;
		words: number;
		estimatedChunks: number;
	}> = [];

	// Process all files
	const allFiles = [
		...docFiles.map((f) => ({ path: f, type: "doc" as const })),
		...blogFiles.map((f) => ({ path: f, type: "blog" as const })),
	];

	for (const file of allFiles) {
		const fullPath = path.join(docsBasePath, file.path);
		const fileContent = fs.readFileSync(fullPath, "utf-8");
		const { content } = matter(fileContent);

		const words = content.split(/\s+/).length;
		const chars = content.length;

		// Estimate chunks based on configuration
		// Assuming ~1.5 characters per token (rough estimate)
		// Max chunk size: 1500 tokens ≈ 2250 characters
		// Overlap: 150 tokens ≈ 225 characters
		const chunkSize = 2250;
		const overlap = 225;
		const effectiveChunkSize = chunkSize - overlap;
		const estimatedChunks = Math.ceil(chars / effectiveChunkSize);

		totalWords += words;
		totalChars += chars;
		fileStats.push({
			path: file.path,
			words,
			estimatedChunks,
		});
	}

	// Calculate totals
	const totalEstimatedChunks = fileStats.reduce(
		(sum, stat) => sum + stat.estimatedChunks,
		0
	);

	console.log("📈 Statistics:");
	console.log(`Total words: ${totalWords.toLocaleString()}`);
	console.log(`Total characters: ${totalChars.toLocaleString()}`);
	console.log(
		`Average words per file: ${Math.round(totalWords / allFiles.length)}`
	);
	console.log(`\nEstimated total chunks: ${totalEstimatedChunks}`);

	// Calculate embedding requirements
	const batchSize = 5;
	const delayPerBatch = 3; // seconds
	const maxChunksPerRun = 100;

	const totalBatches = Math.ceil(totalEstimatedChunks / batchSize);
	const runsNeeded = Math.ceil(totalEstimatedChunks / maxChunksPerRun);
	const timePerRun = Math.ceil(
		((maxChunksPerRun / batchSize) * delayPerBatch) / 60
	); // minutes
	const totalTime = timePerRun * runsNeeded;

	console.log("\n⏱️  Embedding Estimates (with current rate limiting):");
	console.log(`Total batches: ${totalBatches}`);
	console.log(`Runs needed: ${runsNeeded}`);
	console.log(`Time per run: ~${timePerRun} minutes`);
	console.log(
		`Total time: ~${totalTime} minutes (${Math.round(totalTime / 60)} hours)`
	);

	console.log("\n🔢 Rate Limit Usage:");
	console.log(
		`RPM (Requests Per Minute): ${batchSize * 20}/100 per run (max 20 batches/min)`
	);
	console.log(`RPD (Requests Per Day): ${maxChunksPerRun}/1000 per run`);
	console.log(`Total API calls: ${totalEstimatedChunks}`);

	// Show top 10 largest files
	console.log("\n📄 Top 10 largest files (by estimated chunks):");
	const topFiles = fileStats
		.sort((a, b) => b.estimatedChunks - a.estimatedChunks)
		.slice(0, 10);

	topFiles.forEach((file, index) => {
		console.log(`${index + 1}. ${file.path}`);
		console.log(
			`   Words: ${file.words}, Est. chunks: ${file.estimatedChunks}`
		);
	});

	console.log("\n✅ Estimation complete!");
	console.log("\nTo start embedding, run: npm run embed-docs");
}

// Run the estimation
estimateChunks().catch(console.error);
