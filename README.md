# envoy-sync

> CLI tool to sync and diff environment variable sets across multiple `.env` files and remote secret stores.

---

## Installation

```bash
go install github.com/yourusername/envoy-sync@latest
```

Or download a prebuilt binary from the [Releases](https://github.com/yourusername/envoy-sync/releases) page.

---

## Usage

```bash
# Diff two .env files
envoy-sync diff .env.staging .env.production

# Sync variables from a remote store to a local .env file
envoy-sync sync --from aws-secrets://my-app/prod --to .env.local

# Push local .env variables to a remote secret store
envoy-sync push --from .env.production --to aws-secrets://my-app/prod

# List all keys missing between two sources
envoy-sync diff --missing-only .env .env.example
```

### Supported Sources

| Source | Format |
|---|---|
| Local file | `.env`, `.env.*` |
| AWS Secrets Manager | `aws-secrets://<secret-name>` |
| HashiCorp Vault | `vault://<path>` |

---

## Configuration

Optional config file at `.envoy-sync.yaml`:

```yaml
default_store: aws-secrets://my-app/prod
ignore_keys:
  - DEBUG
  - LOG_LEVEL
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE) © 2024 Your Name