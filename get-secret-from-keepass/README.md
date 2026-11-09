# get-secret-from-keepass

Print the password of a KeePass entry by its path.

The path is `/<group>/.../<entry>`, where all parts except the last are
group names (starting with the root group, e.g. `Root`) and the last part
is the entry title.

## Install (master)

```
go install github.com/sikalabs/go-scripts/get-secret-from-keepass@master
```

## Example Usage

```bash
get-secret-from-keepass secrets.kdbx /Root/group1/bar_entry --password changeme
```

or with the password from the env var

```bash
KEEPASS_MASTER_PASSWORD=changeme get-secret-from-keepass secrets.kdbx /Root/group1/bar_entry
```

or with the password from a file (trailing newline is stripped)

```bash
get-secret-from-keepass secrets.kdbx /Root/group1/bar_entry --password-file ./password.txt
```

`--password` and `--password-file` cannot be combined; both take precedence over `KEEPASS_MASTER_PASSWORD`.
