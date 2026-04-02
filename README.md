# go-symphony 🎵

![go-symphony](./assets/go-symphony.png)


Inspired by [go-blueprint](https://github.com/Melkeydev/go-blueprint), I aim to develop a modern CLI tool that generates Go web projects. go-symphony produces Go applications utilizing the Gin framework, PostgreSQL database, and SQLC for type-safe queries. My motivation stems from exclusively employing Gin as an API server and PostgreSQL as my database, eliminating the need for alternative options. Having recently adopted Supabase extensively, I intend to integrate its capabilities into this project. Finally, as I use SvelteKit or Next.js to build the front-end, so I incorporate it into go-symphony. 
## 📺 Demo

Check out a video walkthrough of go-symphony in action: [YouTube Demo](https://www.youtube.com/watch?v=djzhFAd4rjg)
## ✨ Features

- **🚀 Modern Go Architecture**: Clean project structure with best practices built-in
- **🌐 Gin Framework**: High-performance HTTP web framework with httprouter
- **🗄️ PostgreSQL Support**: PostgreSQL with pgx driver for optimal performance
- **🔒 Type-Safe Queries**: SQLC integration for compile-time SQL validation
- **☁️ Supabase Integration**: Modern backend-as-a-service with auth and real-time features
- **🐳 Docker Ready**: Optional Docker and docker-compose configuration
- **⚡ Hot Reload**: Air integration for development with live reload
- **🧪 Testing Setup**: Pre-configured testing with testcontainers
- **🔄 CI/CD**: GitHub Actions workflows for testing and releasing
- **🌐 Frontend Integration**: Optional SvelteKit and Next.js frontend generation
- **🔌 WebSocket Support**: Real-time communication capabilities

## 📦 Installation

### Install with Go

```bash
go install github.com/Tomlord1122/go-symphony@latest
```

### Install from Releases

Download the latest binary from [GitHub Releases](https://github.com/Tomlord1122/go-symphony/releases).

### Build from Source

```bash
git clone https://github.com/Tomlord1122/go-symphony.git
cd go-symphony
go build -o go-symphony .
```

## 📋 Prerequisites

### Required
- **Go 1.23+** to build or run the CLI

### Optional
- **Git** for `--git stage` or `--git commit`
- **Docker** for generated Postgres containers and Docker assets
- **Supabase CLI** for Supabase bootstrap flows
- **Node.js / npx** for SvelteKit or Next.js bootstrap flows
- **pnpm** if you want pnpm-based frontend installs

## 🚀 How To Use

### Interactive Create Flow

Use this when you want prompts:

```bash
go-symphony create
```

Use advanced feature prompts:

```bash
go-symphony create --advanced
```

### Non-Interactive Create Flow

Use this when scripting or working with an AI agent:

```bash
go-symphony create \
  --name github.com/acme/my-api \
  --driver postgres \
  --feature sqlc \
  --feature docker \
  --frontend none \
  --git skip \
  --no-interactive
```

### Preview The Plan Without Writing Files

```bash
go-symphony create \
  --name github.com/acme/my-api \
  --driver postgres \
  --feature sqlc \
  --frontend none \
  --git skip \
  --no-interactive \
  --dry-run
```

### Use The Dedicated Plan Command

`plan` is the agent-friendly preview command. It validates the input and prints the scaffold steps without creating files.

```bash
go-symphony plan \
  --name github.com/acme/my-api \
  --driver postgres \
  --feature sqlc \
  --feature docker \
  --frontend sveltekit \
  --sveltekit-template minimal \
  --sveltekit-types ts \
  --sveltekit-package-manager pnpm \
  --git skip
```

### Output JSON For Automation

Both `create --dry-run` and `plan` support JSON output.

```bash
go-symphony plan \
  --name github.com/acme/my-api \
  --driver supabase \
  --supabase-mode init-only \
  --frontend none \
  --git skip \
  --output json
```

### Skip Install And Formatting Steps

Use `--skip-install` when you only want the scaffolded files and want to run dependency/bootstrap steps later.

```bash
go-symphony create \
  --name github.com/acme/my-api \
  --driver postgres \
  --feature docker \
  --frontend none \
  --git skip \
  --no-interactive \
  --skip-install
```

## 🤖 AI Agent Usage

`go-symphony` now works well in an AI-assisted workflow because it supports:

- `--no-interactive` for deterministic execution
- `--dry-run` for previewing the scaffold plan
- `plan` for a dedicated planning step
- `--output json` for machine-readable output
- `--skip-install` for file generation without dependency/bootstrap side effects

Recommended AI workflow:

1. Generate a plan first.
2. Review the JSON or text output.
3. Run `create` with the same flags.
4. Optionally run follow-up bootstrap steps yourself.

Example:

```bash
go-symphony plan \
  --name github.com/acme/my-api \
  --driver postgres \
  --feature sqlc \
  --feature githubaction \
  --frontend none \
  --git skip \
  --output json

go-symphony create \
  --name github.com/acme/my-api \
  --driver postgres \
  --feature sqlc \
  --feature githubaction \
  --frontend none \
  --git skip \
  --no-interactive \
  --skip-install
```

There is also a reusable `go-symphony` skill definition available in the `tomtom-skill` repository:

- `../tomtom-skill/skills/go-symphony/SKILL.md`

## ⚙️ Main Options

### Database Drivers

- `postgres`
- `supabase`
- `none`

### Advanced Features

- `sqlc`
- `docker`
- `githubaction`
- `websocket`

### Frontend Frameworks

- `sveltekit`
- `nextjs`
- `none`

### Git Modes

- `commit`
- `stage`
- `skip`

### Supabase Modes

- `init-only`
- `local-db`

## 🧱 Generated Project Structure

Typical backend output:

```text
your-project/
├── cmd/api/main.go
├── internal/
│   ├── server/
│   └── database/
├── Makefile
├── README.md
├── .env
├── .air.toml
└── .gitignore
```

Optional additions depend on selected features:

- `sqlc/` for SQLC migrations and queries
- `supabase/` for Supabase bootstrap files
- `.github/workflows/` for GitHub Actions
- `Dockerfile` and `docker-compose.yml` only when `feature=docker` is enabled
- `<project>-frontend/` for SvelteKit or Next.js

## 🛠️ Development Workflow

After generating a project:

```bash
cd your-project

# Start the generated Go server
make run

# Optional: start the Postgres container
make docker-run

# Optional: generate SQLC code
make sqlc-generate

# Optional: run tests
make test

# Optional: build the project
make build
```

If you want generated Docker assets, include `--feature docker` when creating the project.

## 📖 Commands

### Create

```bash
go-symphony create [flags]
```

### Plan

```bash
go-symphony plan [flags]
```

### Version

```bash
go-symphony version
```

## 🗺️ Workflow

![workflow](./assets/go-sym-arc.png)

## 🤝 Contributing

We welcome contributions! Please see our contributing guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments
- Inspired and referenced by [go-blueprint](https://github.com/Melkeydev/go-blueprint)
- Built with [Cobra CLI](https://github.com/spf13/cobra) and [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Claude CLI-inspired color scheme and user experience
