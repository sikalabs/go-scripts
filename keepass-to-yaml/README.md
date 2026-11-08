# keepass-to-yaml

Export a KeePass database to a YAML file. The reverse of
[yaml-to-keepass](../yaml-to-keepass), the output can be fed back to it.

Only title, username, password, URL and notes of each entry are exported.
The recycle bin is skipped.

## Install (master)

```
go install github.com/sikalabs/go-scripts/keepass-to-yaml@master
```

## Example Usage

```bash
keepass-to-yaml secrets.kdbx secrets.yaml --password changeme
```

or with the password from the env var

```bash
KEEPASS_MASTER_PASSWORD=changeme keepass-to-yaml secrets.kdbx secrets.yaml
```

or with the password from a file (trailing newline is stripped)

```bash
keepass-to-yaml secrets.kdbx secrets.yaml --password-file ./password.txt
```

`--password` and `--password-file` cannot be combined; both take precedence over `KEEPASS_MASTER_PASSWORD`.

The output YAML contains the master password and all secrets in plain text
and is written with `0600` permissions.
