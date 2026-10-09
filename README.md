# Go Scripts

Eight beginner-friendly Go DevOps commands, organized as one Go module with a separate command folder for each tool. Each tool has small helper functions and explicit error handling. Go 1.22 or newer is required; all Go code uses the standard library.

| Command | Approach |
| --- | --- |
| [health-check](cmd/health-check/) | Concurrent service health checks |
| [docker-monitor](cmd/docker-monitor/) | Docker container state inspection |
| [log-analyzer](cmd/log-analyzer/) | Regex log analysis |
| [disk-monitor](cmd/disk-monitor/) | Available disk space checks |
| [pod-checker](cmd/pod-checker/) | Kubernetes API pod inspection |
| [server-inventory](cmd/server-inventory/) | Human-readable Linux inventory |
| [deploy-cli](cmd/deploy-cli/) | Deployment webhook CLI |
| [api-agent](cmd/api-agent/) | State-change monitoring agent |

## Start here

```sh
git clone https://github.com/navaparvaz/Go-Scripts.git
cd Go-Scripts
go run ./cmd/health-check http://localhost:8080/health
```

Run commands from this repository root. Read the command README for prerequisites and options. Flags precede positional arguments. Docker needs its installed CLI; Kubernetes checks use a local kubectl proxy; inventory and disk checks require Linux. Deployment uses an HTTP webhook, with preview as the default.

## Structure

- `cmd/<tool>/`: command source and its own README
- `examples/`: sample inputs
- `go.mod`: single module for all commands
- `bin/`: optional compiled binaries, ignored by Git

## Validate

```sh
go test ./...
go vet ./...
go build ./...
```

Tests use local HTTP servers and temporary files; no live deployments or cluster access are needed. These examples print alerts locally. The agent logs state transitions and can run under a process supervisor.
