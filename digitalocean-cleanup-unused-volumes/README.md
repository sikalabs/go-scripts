# digitalocean-cleanup-unused-volumes

Find and delete DigitalOcean block storage volumes not attached to any
droplet (e.g. leftover PVCs from deleted Kubernetes clusters). Runs as a dry run by
default, use `--delete` to actually delete them.

For safety, the script only runs against teams listed in the static
`eligibleTeamUUIDs` allowlist in the source code. The team of the token is
checked first and the script exits for any other team.

## Install (master)

```
go install github.com/sikalabs/go-scripts/digitalocean-cleanup-unused-volumes@master
```

## Example Usage

Dry run (list unused volumes only):

```bash
digitalocean-cleanup-unused-volumes --token dop_v1_xxx
```

Delete unused volumes:

```bash
digitalocean-cleanup-unused-volumes --token dop_v1_xxx --delete
```

## Configuration Sources

The token can be set via `--token`, the `DIGITALOCEAN_TOKEN` env var, or a
`.env` file in the current directory. Precedence: flag > env var > `.env`.

```bash
# .env
DIGITALOCEAN_TOKEN=dop_v1_xxx
```
