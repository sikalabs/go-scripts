# s3-backups-cleanup

Clean up backups in an S3 bucket using this retention policy:

- keep one backup per month forever (the oldest object of each month)
- keep one backup per day for the last 30 days (the newest object of each day)
- delete everything else

Backups are identified by the date in the file name, in the format
`YYYY-MM-DD_HH-MM-SS` (e.g. `sikalabs-com.2026-10-22_18-05-09_mqk.sql.gz`).
The retention is applied separately to each name (the part before the date,
e.g. `sikalabs-com.`), so one bucket can hold backups of multiple sites.
Objects without a date in the name are skipped and never deleted.

Runs as a dry run by default, use `--delete` to actually delete objects.

## Install (master)

```
go install github.com/sikalabs/go-scripts/s3-backups-cleanup@master
```

## Example Usage

Dry run (list what would be kept and deleted):

```bash
s3-backups-cleanup \
  --access-key xxx \
  --secret-key xxx \
  --bucket-name my-backups
```

Delete old backups:

```bash
s3-backups-cleanup \
  --access-key xxx \
  --secret-key xxx \
  --bucket-name my-backups \
  --delete
```

### S3-compatible services (e.g. MinIO, DigitalOcean Spaces)

Path-style addressing is used automatically with a custom endpoint.

```bash
s3-backups-cleanup \
  --access-key xxx \
  --secret-key xxx \
  --endpoint https://fra1.digitaloceanspaces.com \
  --region fra1 \
  --bucket-name my-backups
```

## Flags

```
  -a, --access-key string    AWS Access Key
  -b, --bucket-name string   Bucket Name
      --delete               Actually delete objects (default is dry run)
  -e, --endpoint string      Custom S3 Endpoint
  -r, --region string        AWS Region (default "us-east-1")
  -s, --secret-key string    AWS Secret Key
```
