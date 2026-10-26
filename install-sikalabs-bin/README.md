# install-sikalabs-bin

Install a binary from Sikalabs binary storage (S3 on DigitalOcean). Linux amd64 only.

Downloads `https://sikalabs.fra1.digitaloceanspaces.com/bin/<name>-linux-amd64`
and installs it as `/usr/local/bin/<name>`.

## Install (master)

```
go install github.com/sikalabs/go-scripts/install-sikalabs-bin@master
```

or without Go

```bash
curl -fsSL https://sikalabs.fra1.digitaloceanspaces.com/bin/install-sikalabs-bin-linux-amd64 -o /usr/local/bin/install-sikalabs-bin
chmod +x /usr/local/bin/install-sikalabs-bin
```

## Example Usage

```bash
install-sikalabs-bin ci-build-example
```

or into a custom directory

```bash
install-sikalabs-bin ci-build-example --dir ~/bin
```
