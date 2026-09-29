# digitalocean-cleanup-empty-loadbalancers

Find and delete DigitalOcean load balancers without any droplets attached
(e.g. leftovers from deleted Kubernetes clusters). Runs as a dry run by
default, use `--delete` to actually delete them.

For safety, the script only runs against teams listed in the static
`eligibleTeamUUIDs` allowlist in the source code. The team of the token is
checked first and the script exits for any other team.

## Install (master)

```
go install github.com/sikalabs/go-scripts/digitalocean-cleanup-empty-loadbalancers@master
```

## Example Usage

Dry run (list empty load balancers only):

```bash
digitalocean-cleanup-empty-loadbalancers --token dop_v1_xxx
```

Delete empty load balancers:

```bash
digitalocean-cleanup-empty-loadbalancers --token dop_v1_xxx --delete
```

## Configuration Sources

The token can be set via `--token`, the `DIGITALOCEAN_TOKEN` env var, or a
`.env` file in the current directory. Precedence: flag > env var > `.env`.

```bash
# .env
DIGITALOCEAN_TOKEN=dop_v1_xxx
```
