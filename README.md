# Websocket-Webhook-Relay Installation & Quick Start

This project is a simple relay from a websocket to a webhook, written in Go.

## Requirements

- Go (version 1.18+ recommended)
- Docker (optional, for containerized usage)
- Git

---

## 1. Build & Run with Go

Clone the repository:

```bash
git clone https://github.com/alpha-strike-space/Websocket-Webhook-Relay.git
cd Websocket-Webhook-Relay
```

Build the project if any modifications are made:

```bash
go mod init websocket-relay
```
```bash
go mod tidy
```

> Optionally, check if the project requires configuration (e.g., WEBHOOK_URL). Review the README or source code for details.

---

## 2. Run with Docker

Docker Environment:

Build and run image:

```bash
docker compose up -d
```

Clear container:

```bash
docker compose down
```
---

## Configuration

Most Go web servers accept configuration via environment variables or config files. Check the source code for specifics.

---

## Support

For questions or issues, open an [issue](https://github.com/alpha-strike-space/Websocket-Webhook-Relay/issues) in the repo.
