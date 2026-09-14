# vault-secret-sync

Recursively sync secrets at a given path from one HashiCorp Vault to another.
Supports both KV v1 and KV v2 secrets engines (auto-detected per mount).

## Install (master)

```
go install github.com/sikalabs/go-scripts/vault-secret-sync@master
```

## Example Usage

```bash
vault-secret-sync \
  --source-addr https://vault-source.example.com \
  --source-token hvs.source-token \
  --target-addr https://vault-target.example.com \
  --target-token hvs.target-token \
  --sync-path secret/my-app
```

`--sync-path` must include the mount name (e.g. `secret/my-app`), and is
synced recursively from the source Vault to the same path on the target
Vault.
