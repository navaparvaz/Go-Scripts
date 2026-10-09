# Kubernetes API pod inspection

First run kubectl proxy --port=8001 in another terminal with your configured kubeconfig. This tool calls the HTTP API through that local proxy and follows paginated responses. Checks Ready condition and total container restart count per pod. Succeeded job pods pass readiness but still use the restart limit. An empty namespace returns success with zero pods. -api changes the proxy URL; direct kubeconfig, bearer-token and custom-CA loading are not implemented. [Kubernetes proxy documentation](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_proxy/).

## Run from the repository root

```sh
go run ./cmd/pod-checker -namespace default -max-restarts 5
```

## Build

```sh
mkdir -p bin
go build -o bin/pod-checker ./cmd/pod-checker
```

Read `main.go` starting at `run()`, then follow its small helper functions. `run()` returns an exit code so resources can close before `main()` exits.

Exit codes: 0 = completed successfully, 1 = unhealthy check or operational failure, 2 = invalid usage. Log matches produce 1; the agent continues on unhealthy results.
