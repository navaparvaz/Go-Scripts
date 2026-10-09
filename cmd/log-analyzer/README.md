# Regex log analysis

Accepts multiple log file arguments, prints matching lines, and reports per-file and combined counts. Default pattern matches whole words error/fatal/panic without regard to case. Set -pattern for a custom Go regular expression. Maximum line size is 1 MiB.

## Run from the repository root

```sh
go run ./cmd/log-analyzer examples/sample.log
```

## Build

```sh
mkdir -p bin
go build -o bin/log-analyzer ./cmd/log-analyzer
```

Read `main.go` starting at `run()`, then follow its small helper functions. `run()` returns an exit code so resources can close before `main()` exits.

Exit codes: 0 = completed successfully, 1 = unhealthy check or operational failure, 2 = invalid usage. Log matches produce 1; the agent continues on unhealthy results.
