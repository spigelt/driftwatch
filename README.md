# driftwatch

Lightweight daemon that detects configuration drift between running containers and their declared compose definitions.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your Compose file and let it watch for drift:

```bash
driftwatch --compose docker-compose.yml --interval 30s
```

Example output when drift is detected:

```
[DRIFT] web: image mismatch — expected nginx:1.25, got nginx:1.23
[DRIFT] api: environment variable PORT missing from running container
[OK]    db: no drift detected
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--compose` | `docker-compose.yml` | Path to Compose file |
| `--interval` | `60s` | How often to check for drift |
| `--notify` | `""` | Webhook URL for drift alerts |
| `--once` | `false` | Run a single check and exit |

Run a one-shot check and exit with a non-zero code if drift is found:

```bash
driftwatch --compose docker-compose.yml --once
```

---

## How It Works

`driftwatch` reads your Compose file, inspects the currently running containers via the Docker API, and compares image tags, environment variables, port bindings, and volume mounts. Any discrepancy is reported as drift.

---

## License

MIT © 2024 yourusername