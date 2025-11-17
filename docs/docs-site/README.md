# go-infra Documentation Site

Documentation website for go-infra, built with [Docusaurus](https://docusaurus.io/).

## 🚀 Quick Start

### Installation

```bash
npm install
```

### Local Development

```bash
npm start
```

This command starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

### Build

```bash
npm run build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

### Deployment

```bash
npm run deploy
```

## 📁 Documentation Structure

```
docs/
├── intro.md                          # Introduction
├── getting-started/
│   ├── installation.md               # Installation guide
│   └── quick-start.md                # Quick start tutorial
├── core-concepts/
│   └── architecture.md               # Architecture overview
├── http-server/
│   ├── building-apis.md              # Building REST APIs
│   └── crud-operations.md            # CRUD operations
├── database/
│   └── setup.md                      # Database setup
├── authentication/
│   └── getting-started.md            # Authentication guide
├── components/
│   ├── logger/                       # Logger documentation
│   │   ├── index.md                  # Logger overview
│   │   ├── quickstart.md             # Quick start guide
│   │   ├── configuration.md          # Configuration options
│   │   ├── usage.md                  # Usage guide
│   │   ├── fx-integration.md         # Fx integration
│   │   ├── adapters.md               # External adapters
│   │   └── best-practices.md         # Best practices
│   └── crypto.md                     # Crypto documentation
└── examples/
    └── users-api.md                  # Users API example
```

## ✍️ Writing Documentation

### Creating New Documents

1. Create a new `.md` or `.mdx` file in the appropriate directory
2. Add frontmatter:

```md
---
sidebar_position: 1
---

# Your Title

Your content here...
```

3. The document will automatically appear in the sidebar

### Code Blocks

Use language-specific code blocks:

\`\`\`go
package main

func main() {
fmt.Println("Hello, World!")
}
\`\`\`

### Admonitions

```md
:::tip
This is a helpful tip!
:::

:::info
This is important information.
:::

:::warning
This is a warning.
:::
```

## 🎨 Customization

- **Configuration**: Edit `docusaurus.config.ts`
- **Styling**: Edit `src/css/custom.css`
- **Sidebar**: Edit `sidebars.ts`

## 📦 AI Integration (Planned)

To add AI-powered question answering:

- Follow: https://github.com/ahelmy/docusaurus-ai

## 🔗 Links

- **Live Site**: [http://localhost:3000](http://localhost:3000)
- **go-infra Repository**: https://github.com/phatnt199/go-infra

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally with `npm start`
5. Submit a pull request

## 📄 License

MIT License
