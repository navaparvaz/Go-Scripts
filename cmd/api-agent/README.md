# State-change monitoring agent

Monitors multiple URLs sequentially, pausing after each sweep. Emits JSON lines for each first result and only subsequent healthy/unhealthy transitions. Ctrl+C or SIGTERM cancels an active request and stops the agent. Unhealthy endpoints do not stop monitoring. On Linux/macOS build it, then run ./bin/api-agent -interval 30s http://localhost:8080/health > agent.log 2>&1 &. Manage it with your process supervisor if you need it to survive sessions.

## Run from the repository root

```sh
go run ./cmd/api-agent -interval 30s http://localhost:8080/health http://localhost:3000/health
```

## Build

```sh
mkdir -p bin
go build -o bin/api-agent ./cmd/api-agent
```

Read `main.go` starting at `run()`, then follow its small helper functions. `run()` returns an exit code so resources can close before `main()` exits.

Exit codes: 0 = completed successfully, 1 = unhealthy check or operational failure, 2 = invalid usage. Log matches produce 1; the agent continues on unhealthy results.
