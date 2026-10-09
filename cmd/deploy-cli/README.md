# Deployment webhook CLI

Creates a JSON deployment request. Default mode prints a preview without network access. To send it, add -execute -endpoint http://localhost:8080/deploy. Optional DEPLOY_TOKEN sets a bearer token header (use HTTPS for real credentials). Requires your own webhook that accepts application/version/environment and performs deployment; the tool does not install or provide that server. Redirects are refused; requests are not retried to avoid duplicate deployments. A 2xx response means the request was accepted, not that rollout completed.

## Run from the repository root

```sh
go run ./cmd/deploy-cli deploy -app web -version v1.0 -env staging
```

## Build

```sh
mkdir -p bin
go build -o bin/deploy-cli ./cmd/deploy-cli
```

Read `main.go` starting at `run()`, then follow its small helper functions. `run()` returns an exit code so resources can close before `main()` exits.

Exit codes: 0 = completed successfully, 1 = unhealthy check or operational failure, 2 = invalid usage. Log matches produce 1; the agent continues on unhealthy results.
