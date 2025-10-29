# Releasing a New Version of go-infra

This guide outlines the development and release process for `go-infra`. Following these steps ensures consistency and helps consumers of the module.

## 1. Development Workflow

All new features, bug fixes, and changes should be developed in a feature branch.

### Step 1: Create a Branch

Create a new branch from the `main` branch. Use a descriptive name, like:

- `feature/add-new-component`
- `fix/resolve-memory-leak`
- `docs/update-readme`

```bash
# Make sure your main branch is up-to-date
git checkout main
git pull origin main

# Create your new branch
git checkout -b feature/your-descriptive-name
```

### Step 2: Commit and Push Changes

Make your code changes, then commit them with a clear message.

```bash
git add .
git commit -m "feat: add awesome new feature"
git push origin feature/your-descriptive-name
```

### Step 3: Open a Pull Request

Go to the GitHub repository and open a pull request from your feature branch to the `main` branch. Ensure the PR description clearly explains the changes. Once approved and merged, your changes are ready for the next release.

## 2. Release Workflow

Releasing a new version involves creating a semantic version tag and pushing it to the repository.

### Step 1: Decide the Version Number

This project follows [Semantic Versioning (SemVer)](https://semver.org/). After merging changes into `main`, decide the next version number based on the changes:

- **MAJOR** (`vX.0.0`): For incompatible API changes.
- **MINOR** (`v1.X.0`): For adding functionality in a backward-compatible manner.
- **PATCH** (`v1.0.X`): For backward-compatible bug fixes.

For a first release, `v0.1.0` (initial development) or `v1.0.0` (stable) are good choices.

### Step 2: Create an Annotated Git Tag

From the `main` branch, create an annotated tag. This is critical for the Go module system.

```bash
# Ensure you are on the main branch with the latest changes
git checkout main
git pull origin main

# Create the annotated tag (e.g., v1.1.0)
git tag -a v1.1.0 -m "Release v1.1.0"
```

### Step 3: Push the Tag

Push the new tag to the remote repository. This makes the new version available to consumers.

```bash
git push origin v1.1.0
```

Consumers can now get the new version by running:
`go get github.com/phatnt199/go-infra@v1.1.0`

## Special Rule: Major Versions (v2 and beyond)

When you release a `v2` or higher, Go modules require a special change. The module path in `go.mod` **must** be updated to include the major version suffix (e.g., `/v2`).

### Example: Releasing v2.0.0

1.  **Update `go.mod`:**
    Change `module github.com/phatnt199/go-infra` to:

    ```go
    module github.com/phatnt199/go-infra/v2
    ```

2.  **Update Import Paths:**
    Update all import paths within the project to include `/v2`.

    ```go
    // Old
    import "github.com/phatnt199/go-infra/pkg/logger"
    // New
    import "github.com/phatnt199/go-infra/v2/pkg/logger"
    ```

3.  **Tidy and Test:**

    ```bash
    go mod tidy
    go test ./...
    ```

4.  **Commit, Tag, and Push:**
    Commit the changes, then tag and push `v2.0.0`.
    ```bash
    git commit -m "feat: release v2.0.0"
    git tag -a v2.0.0 -m "Release v2.0.0"
    git push origin v2.0.0
    ```
