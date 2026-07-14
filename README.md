<div align="center">
  <img src="docs/assets/logo.svg" alt="Sooqara" width="160" />

  # Sooqara

  **E-commerce Listing Factory Powered by Agnes AI**

  [![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-c9a84c.svg)](LICENSE)
  [![Go](https://img.shields.io/badge/Go-1.25-1c2f52.svg)](go.mod)
  [![SQLite](https://img.shields.io/badge/DB-SQLite-c9a84c.svg)](internal/store)
  [![HTMX](https://img.shields.io/badge/UI-HTMX-1c2f52.svg)](internal/web)

</div>

---

Sooqara is an **e-commerce listing factory** that automates the end-to-end creation of marketplace listings. Given a product image, it performs vision analysis, generates compliant copy, produces multiple image variants, creates async video assets, and exports everything in platform-ready formats — all orchestrated through a deterministic pipeline with rate limiting, cost tracking, and full audit trails.

The pipeline is built around one idea: **automated listings need deterministic guardrails.** Every stage — from vision analysis through export — is bounded by token budgets, content constraints, and validation gates. No stage proceeds without passing its checks.

---

## How it works

```
  UPLOAD              VISION ANALYSIS           COPY GENERATION        IMAGE VARIANTS        VIDEO CREATION        EXPORT
  product image      → Agnes AI vision         → constraint-bound     → seed-locked         → async poll          → ZIP / CSV / TSV
                       → prompt extraction      → banned-phrase filter → multi-variant gen  → progress SSE        → JSON
                       → style enforcement      → tone/style match     → blob storage       → cost tracking       → platform-ready
```

1. **Upload & ingest.** A product image is uploaded through the HTMX web interface or REST API. The file is stored in the local blob store and a job record is created in SQLite with `status: queued`.

2. **Vision analysis.** The Agnes AI vision model analyzes the image, extracts product attributes, suggests a title and description, identifies dominant colors and style cues, and generates a structured prompt for downstream stages.

3. **Copy generation.** A chat model produces marketplace-ready copy — title, bullet points, description — constrained by character limits, banned-phrase filters, and the style profile derived from vision analysis.

4. **Image variants.** Seed-locked image generation creates multiple stylistic variants of the product image. Each variant is stored in the blob store with its seed, prompt, and generation cost tracked.

5. **Async video.** An image-to-video model generates a short product showcase video from the primary image. Creation is asynchronous; progress is streamed via Server-Sent Events to the web UI.

6. **Export.** Completed listings are packaged into platform-ready formats: ZIP archive, Shopify CSV, WooCommerce CSV, Amazon TSV, or structured JSON. Each export includes all artifacts and metadata.

<div align="center">
  <img src="docs/assets/architecture-diagram.png" alt="Sooqara architecture: upload, vision analysis, copy generation, image variants, async video, export" width="100%" />
</div>

---

## Key properties

| Property | Guarantee |
|---|---|
| **Deterministic pipeline** | Every listing follows a fixed six-stage workflow — vision, copy, images, video, export — with no shortcuts or reordering. |
| **Rate-limited AI calls** | Token-bucket limiter with deterministic `fakeClock` for testing; RPM, concurrency, and budget are configurable via `.env`. |
| **Seed-locked variants** | Image generation uses locked seeds for reproducibility; each variant's seed, prompt, and output are stored and queryable. |
| **Cost tracking** | Per-call cost ledger with SQLite persistence, token accounting, and configurable pricing model for chat, image, and video. |
| **Activity journal** | Immutable, append-only SQLite log of every pipeline event — job creation, stage transitions, errors, cost entries. |
| **Hardened input** | Title sanitization, URL validation, API key length checks, tag limits, path-traversal prevention, header allowlists. |
| **Multi-format export** | ZIP bundles, Shopify CSV, WooCommerce CSV, Amazon TSV, and machine-readable JSON — all with slugified columns and proper escaping. |
| **Real-time UI** | HTMX-driven server-rendered interface with SSE progress streams, no build step, no JavaScript framework dependencies. |

---

## Platform capabilities

Sooqara is a modular Go monolith organized across twelve internal packages. The table below summarizes each capability; the package paths link to the authoritative source.

| Capability | Package | Description |
|---|---|---|
| **Configuration** | [`internal/config`](internal/config) | Typed `Config` struct with environment variable loading, defaults, and `Validate()` with joined errors. |
| **Rate Limiter** | [`internal/limiter`](internal/limiter) | Token-bucket limiter with `fakeClock` for deterministic concurrency testing and budget enforcement. |
| **Activity Journal** | [`internal/journal`](internal/journal) | Append-only SQLite log recording every pipeline event with query support by job and stage. |
| **AI Provider** | [`internal/provider`](internal/provider) | Abstraction over Agnes AI APIs (chat, image, video) with a fake implementation for testing. |
| **Storage** | [`internal/store`](internal/store) | Job and artifact CRUD with SQLite migrations, state transitions, and filesystem blob storage. |
| **Pipeline Stages** | [`internal/pipeline`](internal/pipeline) | Vision analysis, copy generation, image variants, and video creation stages. |
| **Orchestrator** | [`internal/orchestrator`](internal/orchestrator) | Worker pool that dispatches jobs through the pipeline with concurrency control. |
| **HTTP API** | [`internal/api`](internal/api) | REST endpoints for job management, progress SSE, and health checks. |
| **Web UI** | [`internal/web`](internal/web) | HTMX-driven interface with embedded templates, static assets, and real-time progress. |
| **Export** | [`internal/export`](internal/export) | ZIP, Shopify CSV, WooCommerce CSV, Amazon TSV, and JSON writers with slugification. |
| **Observability** | [`internal/observability`](internal/observability) | Metrics counters, cost ledger, pricing model, and validation. |
| **Hardening** | [`internal/hardening`](internal/hardening) | Input sanitization, URL validation, API key checks, tag limits, path traversal prevention. |
| **Release** | [`internal/release`](internal/release) | Build flag injection, manifest I/O, tarball packaging, checksum generation, environment gate. |
| **Version** | [`internal/version`](internal/version) | Build-time version, commit hash, and build timestamp via `-ldflags`. |

---

## Tech stack

| Layer | Stack |
|---|---|
| Language | Go 1.25 |
| Database | SQLite (`modernc.org/sqlite`) — jobs, artifacts, journal, cost ledger |
| AI Integration | Agnes AI API — vision, chat, image generation, image-to-video |
| Web UI | HTMX + embedded Go templates — no JavaScript framework, no build step |
| Export | ZIP archive, Shopify CSV, WooCommerce CSV, Amazon TSV, JSON |
| Observability | Cumulative metrics, per-call cost ledger, configurable pricing model |
| Release | `-ldflags` version injection, SHA-256 checksums, tarball packaging, environment gate |

---

## Development workflow

Sooqara was built using a **branch-per-phase, atomic-commits** workflow:

| Phase | Branch | Description |
|---|---|---|
| 00 | `phase-00` | Repository bootstrap: `go.mod`, `.gitignore`, `Makefile`, `README`, `LICENSE`, typed config |
| 01 | `phase-01` | Rate limiter (token bucket + `fakeClock`) and SQLite activity journal |
| 02 | `phase-02` | Agnes AI provider client (chat, image, video) with fake stub |
| 03 | `phase-03` | Storage layer: job/artifact CRUD, blob storage, SQLite migrations |
| 04 | `phase-04` | Vision analysis stage: image analysis, prompt generation, style enforcement |
| 05 | `phase-05` | Copy generation stage: constraint enforcement, banned phrases |
| 06 | `phase-06` | Seed-locked image variants + async video creation |
| 07 | `phase-07` | Orchestrator and worker pool |
| 08 | `phase-08` | HTTP API with SSE progress events |
| 09 | `phase-09` | HTMX web UI with embedded templates and static assets |
| 10 | `phase-10` | Export pipeline: ZIP, Shopify CSV, WooCommerce CSV, Amazon TSV, JSON |
| 11 | `phase-11` | Observability: metrics counters, cost ledger, pricing model |
| 12 | `phase-12` | Hardening: input validation, sanitization, path traversal prevention |
| 13 | `phase-13` | Packaging: version injection, build flags, tarball packing, checksums |
| 14 | `phase-14` | Final release: manifest I/O, environment gate, full pipeline integration |

Each phase was developed on its own branch, pushed to the remote, and merged to `main` via pull request. All phases pass `go vet` and `go test -race`.

---

## Monorepo layout

```
sooqara/
├── cmd/sooqara/        # Application entry point — loads config, starts server
├── internal/           # Private packages — config, limiter, journal, provider, store
│   ├── api/            # HTTP handlers — job management, SSE, health check
│   ├── config/         # Typed configuration with env loading and validation
│   ├── export/         # Multi-format exporters — ZIP, CSV, TSV, JSON
│   ├── hardening/      # Input validation — title, URL, API key, tags, paths
│   ├── journal/        # Append-only activity log — SQLite-backed
│   ├── limiter/        # Token-bucket rate limiter — with fakeClock for tests
│   ├── observability/  # Metrics, cost ledger, pricing model
│   ├── orchestrator/   # Worker pool — job dispatch and stage orchestration
│   ├── pipeline/       # Vision, copy, image, video pipeline stages
│   ├── provider/       # Agnes AI abstraction — chat, image, video + fake impl
│   ├── release/        # Build flags, manifest I/O, tarball, checksums
│   ├── store/          # Job/artifact persistence — SQLite + blob storage
│   ├── version/        # Build-time version, commit, timestamp
│   └── web/            # HTMX templates, static assets, embedded server
├── docs/               # Assets, diagrams, and reference documentation
├── temp/               # Scratch space — not shipped in production builds
├── Makefile            # build, run, test, fmt, vet, lint, clean, docker
├── .env.example        # Environment variable defaults
└── go.mod              # Go module definition
```

---

## Getting started

### Prerequisites

- Go 1.25+
- SQLite (bundled via `modernc.org/sqlite` — no system dependency)
- An Agnes AI API key

### Quickstart

```bash
# Clone and enter the project
git clone https://github.com/YASSERRMD/Sooqara.git
cd Sooqara

# Copy and configure environment
cp .env.example .env
# Edit .env and set your AGNES_API_KEY

# Build and run
make build
./bin/sooqara
```

The server starts on `:8080` by default. Open `http://localhost:8080` in your browser to access the HTMX web interface.

### Common commands

```bash
make build        # Build the binary to bin/sooqara
make run          # Build and run (requires .env)
make test         # Run all tests with race detection
make fmt          # Format all Go source files
make vet          # Run go vet across all packages
make lint         # Alias for vet
make clean        # Remove bin/ and *.db files
make docker       # Build Docker image
```

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `AGNES_API_KEY` | *(required)* | Agnes AI API authentication key |
| `AGNES_BASE_URL` | `https://apihub.agnes-ai.com/v1` | Agnes AI API base URL |
| `AGNES_POLL_URL` | `https://apihub.agnes-ai.com/agnesapi` | Agnes AI polling endpoint |
| `SOOQARA_RPM` | `18` | Maximum API calls per minute |
| `SOOQARA_DB` | `./sooqara.db` | SQLite database path |
| `SOOQARA_STORAGE` | `./storage` | Blob storage directory |
| `SOOQARA_ADDR` | `:8080` | HTTP server listen address |
| `SOOQARA_WORKERS` | `3` | Number of orchestrator workers |
| `SOOQARA_LOG_LEVEL` | `info` | Logging level (`debug`, `info`, `warn`, `error`) |

---

## Building releases

```bash
# Development build
make build

# Production release with version injection
make release-version VERSION=v1.0.0 COMMIT=$(git rev-parse --short HEAD)

# Or use environment variables
VERSION=v1.0.0 COMMIT=abc123 make release-version
```

The release system generates:
- A versioned binary with embedded build metadata
- SHA-256 checksums for all artifacts
- A JSON manifest describing the release
- A tarball containing all artifacts

---

## Testing

```bash
# Run all tests with race detector
make test

# Run a specific package
go test -race -count=1 ./internal/limiter/

# View coverage
go test -race -count=1 -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Every package includes comprehensive tests covering:
- Unit tests for all public functions and methods
- Concurrency safety tests with `go test -race`
- Compile-time interface verification
- Integration tests for full pipeline flows
- Edge cases and error paths

---

## Architecture highlights

### Rate limiter with deterministic testing

The token-bucket limiter uses a `fakeClock` interface that allows tests to advance time deterministically, verifying that concurrent goroutines correctly compete for tokens without actual wall-clock delays.

### Activity journal

Every pipeline event is appended to an immutable SQLite journal. Queries can filter by job ID, stage, or time range. The journal supports audit trails and cost reconciliation.

### Seed-locked image generation

Image variants are generated with fixed seeds for reproducibility. Each variant's seed, prompt, and output path are stored alongside the artifact in the SQLite database.

### Multi-format export

Exports use slugified column headers, proper CSV quoting, and platform-specific formatting. The ZIP exporter bundles all artifacts with a manifest, while format-specific writers handle Shopify, WooCommerce, and Amazon requirements.

---

## License

Apache 2.0 — see [`LICENSE`](LICENSE).
