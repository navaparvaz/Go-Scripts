# Docker container state inspection

Requires Docker CLI and daemon. With names supplied, inspects those containers; without names, discovers all containers including stopped ones. Uses structured inspect data. A running container with no health check is accepted; a configured health check must be healthy (starting is reported as a failed check). [Docker inspect documentation](https://docs.docker.com/reference/cli/docker/container/inspect/).

## Run from the repository root

```sh
go run ./cmd/docker-monitor web database
```

## Build

```sh
mkdir -p bin
go build -o bin/docker-monitor ./cmd/docker-monitor
```

Read `main.go` starting at `run()`, then follow its small helper functions. `run()` returns an exit code so resources can close before `main()` exits.

Exit codes: 0 = completed successfully, 1 = unhealthy check or operational failure, 2 = invalid usage. Log matches produce 1; the agent continues on unhealthy results.
