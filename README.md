# envlens

> Diff and audit environment variable sets across multiple `.env` files and deployment configs.

---

## Installation

```bash
go install github.com/yourname/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envlens.git && cd envlens && go build -o envlens .
```

---

## Usage

Compare two `.env` files and see what's missing, added, or changed:

```bash
envlens diff .env.development .env.production
```

Audit multiple configs against a baseline:

```bash
envlens audit --baseline .env.example .env.staging .env.production
```

Example output:

```
[MISSING]  .env.production  →  DATABASE_POOL_SIZE
[EXTRA]    .env.staging     →  DEBUG_VERBOSE
[CHANGED]  .env.production  →  LOG_LEVEL  (info → warn)
```

### Flags

| Flag | Description |
|------|-------------|
| `--baseline` | File to treat as the source of truth |
| `--format` | Output format: `text` (default), `json`, `csv` |
| `--ignore` | Comma-separated list of keys to skip |
| `--strict` | Exit with non-zero code if any diff is found |

---

## Why envlens?

Misconfigured environment variables are a common source of bugs and security issues across environments. `envlens` gives you a fast, scriptable way to keep your configs consistent and auditable — locally or in CI pipelines.

---

## License

MIT © [yourname](https://github.com/yourname)