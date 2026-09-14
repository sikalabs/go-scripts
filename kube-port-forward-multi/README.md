# kube-port-forward-multi

Run multiple `kubectl port-forward` sessions at once, each defined as
`<local-port>:<svc/name|pod/name>:<remote-port>` (current namespace) or
`<local-port>:<namespace/svc/name|namespace/pod/name>:<remote-port>` (specific
namespace).

## Install (master)

```
go install github.com/sikalabs/go-scripts/kube-port-forward-multi@master
```

## Example Usage

```bash
kube-port-forward-multi \
  8080:svc/my-service:80 \
  9090:other-namespace/pod/my-pod:9090
```

Each forward is retried automatically if it drops. Press Ctrl+C to stop all
of them.

## Flags

- `--context` - kubeconfig context to use
- `--kubeconfig` - path to kubeconfig file
