# Available disk space checks

Requires Linux. Reads filesystem statistics directly with Statfs and checks available space in GiB against a minimum, rather than a used percentage. Accepts multiple paths. Reserved administrator blocks are excluded; paths on the same filesystem may report the same space.

## Run from the repository root

```sh
go run ./cmd/disk-monitor -min-free 2 / /tmp
```

## Build

```sh
mkdir -p bin
go build -o bin/disk-monitor ./cmd/disk-monitor
```

Read `main.go` starting at `run()`, then follow its small helper functions. `run()` returns an exit code so resources can close before `main()` exits.

Exit codes: 0 = completed successfully, 1 = unhealthy check or operational failure, 2 = invalid usage. Log matches produce 1; the agent continues on unhealthy results.
