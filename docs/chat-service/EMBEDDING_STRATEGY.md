# Embedding Strategy for Free Tier Google Gemini API

## Problem

The original embedding strategy was hitting Google's free tier API limits:

- **RPM (Requests Per Minute)**: 100
- **TPM (Tokens Per Minute)**: 30,000
- **RPD (Requests Per Day)**: 1,000

With 66 documentation files (~25K words), the original strategy would create hundreds of small chunks and process them too quickly, exceeding the rate limits.

## Optimized Solution

### 1. **Larger Chunk Sizes**

- **Before**: 800 tokens per chunk, 100 token overlap
- **After**: 1,500 tokens per chunk, 150 token overlap
- **Impact**: Reduces total chunk count by ~40-50%, fewer API requests needed

### 2. **Smaller Batch Sizes**

- **Before**: 10 chunks per batch
- **After**: 5 chunks per batch
- **Impact**: Stays well under the 100 RPM limit

### 3. **Longer Delays Between Batches**

- **Before**: 1 second delay
- **After**: 3 seconds delay
- **Impact**: Ensures we don't exceed TPM limits (30K tokens/min)

### 4. **Daily Request Limit**

- **New**: Process maximum 100 chunks per run
- **Impact**: Stays well under the 1,000 RPD limit
- **Result**: Script needs to be run multiple times for full embedding

### 5. **Progress Tracking & Resume Capability**

- **New**: Saves progress after each successful batch
- **Location**: `data/embed-progress.json`
- **Benefit**: Can resume from where it left off if interrupted or rate limited
- **Cleanup**: Automatically removes progress file when complete

### 6. **Enhanced Error Handling**

- Detects rate limit errors (429 responses)
- Automatically waits 60 seconds before retrying
- Saves progress before throwing errors
- Provides clear status updates

## Usage

### First Run

```bash
npm run embed-docs
```

This will:

1. Process up to 100 chunks
2. Save progress to `data/embed-progress.json`
3. Report how many chunks remain

### Subsequent Runs

Simply run the same command again:

```bash
npm run embed-docs
```

The script will:

1. Load previous progress
2. Skip already-processed chunks
3. Process the next 100 chunks
4. Update progress

### Monitor Progress

The script outputs:

```
Found 150 already processed. Processing 100 remaining chunks...
Processing batch 1/20 (5 chunks)...
  ✓ Batch completed in 1234ms
  ⏳ Waiting 3000ms to respect rate limits...
```

### Completion

When all chunks are processed:

```
🎉 All documents successfully embedded!
✓ Progress file cleaned up
```

## Estimated Time to Complete

For 66 files (~25K words):

- **Estimated chunks**: ~250-300 (with 1500 token chunks)
- **Chunks per run**: 100
- **Runs needed**: 3-4 runs
- **Time per run**: ~10-15 minutes (with delays)
- **Total time**: 30-60 minutes (spread across multiple runs)

## Tips for Free Tier

1. **Run during off-peak hours**: Less likely to hit rate limits
2. **Run multiple times**: Don't expect one run to complete everything
3. **Be patient**: The delays are necessary to stay within limits
4. **Monitor output**: Check for rate limit warnings
5. **Don't clear progress**: Let the script resume where it left off

## Emergency Recovery

If you need to start over:

```bash
rm data/embed-progress.json
npm run embed-docs
```

## Upgrading to Paid Tier

If you upgrade to a paid Google API plan, you can increase:

- `batchSize`: from 5 to 20+
- `delayBetweenBatches`: from 3000ms to 1000ms
- `maxChunksPerRun`: from 100 to unlimited (remove check)

Edit these values in `src/scripts/embed-docs.ts` in the `embedDocuments` function.

## Alternative: Use Different Embedding Model

Consider switching to:

- **OpenAI embeddings** (better rate limits on paid tier)
- **Local embeddings** (no rate limits, but requires more setup)
- **Cohere embeddings** (generous free tier)

## Database Configuration

The script uses LibSQL (SQLite-compatible) vector database:

- **Location**: `data/docs-chat.db`
- **Dimensions**: 768 (optimal balance for Gemini)
- **Metric**: Cosine similarity

To use Turso (remote database):

```bash
export DATABASE_URL="libsql://your-db-url.turso.io"
export DATABASE_AUTH_TOKEN="your-token"
```
