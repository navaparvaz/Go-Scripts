# Human-readable Linux inventory

Requires Linux /proc and /etc/os-release. Parses key/value files with a reusable helper, including the distribution name. Reports hostname, OS, architecture, logical CPU count, CPU model and total host RAM. Container limits are not substituted for host RAM.

## Run from the repository root

```sh
go run ./cmd/server-inventory
```

## Build

```sh
mkdir -p bin
go build -o bin/server-inventory ./cmd/server-inventory
```

Read `main.go` starting at `run()`, then follow its small helper functions. `run()` returns an exit code so resources can close before `main()` exits.

Exit codes: 0 = completed successfully, 1 = unhealthy check or operational failure, 2 = invalid usage. Log matches produce 1; the agent continues on unhealthy results.
