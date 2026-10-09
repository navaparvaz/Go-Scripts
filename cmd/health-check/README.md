# Concurrent service health checks

Checks up to 50 URLs concurrently using goroutines and a channel. A 2xx status is healthy. Results appear in completion order and include duration. Flags must come before URLs.

## Run from the repository root

```sh
go run ./cmd/health-check http://localhost:8080/health http://localhost:3000/health
```

## Build

```sh
mkdir -p bin
go build -o bin/health-check ./cmd/health-check
```

Read `main.go` starting at `run()`, then follow its small helper functions. `run()` returns an exit code so resources can close before `main()` exits.

Exit codes: 0 = completed successfully, 1 = unhealthy check or operational failure, 2 = invalid usage. Log matches produce 1; the agent continues on unhealthy results.
