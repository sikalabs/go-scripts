# kube-port-forward-multi

Run multiple `kubectl port-forward` sessions at once, each defined as
`<local-port>:<svc/name|pod/name>:<remote-port>`.

## Install (master)

```
go install github.com/sikalabs/go-scripts/kube-port-forward-multi@master
```

## Example Usage

```bash
kube-port-forward-multi \
  8080:svc/my-service:80 \
  9090:pod/my-pod:9090 \
  -n my-namespace
```

Each forward is retried automatically if it drops. Press Ctrl+C to stop all
of them.

## Flags

- `-n, --namespace` - Kubernetes namespace
- `--context` - kubeconfig context to use
- `--kubeconfig` - path to kubeconfig file
