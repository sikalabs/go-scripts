# yaml-to-keepass

Create a KeePass (KDBX 4) database from a YAML file.

## Install (master)

```
go install github.com/sikalabs/go-scripts/yaml-to-keepass@master
```

## Example Usage

```bash
yaml-to-keepass secrets.yaml secrets.kdbx
```

## Input Format

```yaml
master_password: changeme
entries:
  - title: Example
    username: user
    password: secret
    url: https://example.com
    notes: Some notes
groups:
  - name: Servers
    entries:
      - title: prod
        username: root
        password: secret
    groups:
      - name: Nested
        entries: []
```

`master_password` can be omitted from the YAML and set via the
`KEEPASS_MASTER_PASSWORD` env var instead.
