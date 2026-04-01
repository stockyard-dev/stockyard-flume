# Stockyard Flume

**Log aggregator — ship logs via HTTP, store, search, and tail them live**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9210:9210 -v flume_data:/data ghcr.io/stockyard-dev/stockyard-flume
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9210` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9210` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `FLUME_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 2 sources, 50k logs retained | Unlimited sources, 10M logs |
| Price | Free | $2.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Operations & Teams

## License

Apache 2.0
