# Stockyard Crossroads

**Link-in-bio page builder — your domain, your data, no Linktree account**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9300:9300 -v crossroads_data:/data ghcr.io/stockyard-dev/stockyard-crossroads
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9300` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9300` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `CROSSROADS_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 1 profile, 10 links | Unlimited profiles and links |
| Price | Free | $2.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Creator & Small Business

## License

Apache 2.0
