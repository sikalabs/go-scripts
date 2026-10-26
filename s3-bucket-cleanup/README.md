# s3-bucket-cleanup

Delete every object in an S3 bucket. The bucket itself is kept.

Runs as a dry run by default, use `--delete` to actually delete objects.
Before deleting, the script asks for confirmation (like `terraform apply`),
only `yes` is accepted. Use `--yes` to skip it (e.g. in cron), the script
then waits 10 seconds before deleting, so you can still cancel it with Ctrl-C.

## Install (master)

```
go install github.com/sikalabs/go-scripts/s3-bucket-cleanup@master
```

## Example Usage

Dry run (list objects which would be deleted):

```bash
s3-bucket-cleanup \
  --access-key xxx \
  --secret-key xxx \
  --bucket-name my-bucket
```

Delete all objects:

```bash
s3-bucket-cleanup \
  --access-key xxx \
  --secret-key xxx \
  --bucket-name my-bucket \
  --delete
```

### S3-compatible services (e.g. MinIO, DigitalOcean Spaces)

Path-style addressing is used automatically with a custom endpoint.

```bash
s3-bucket-cleanup \
  --access-key xxx \
  --secret-key xxx \
  --endpoint https://fra1.digitaloceanspaces.com \
  --region fra1 \
  --bucket-name my-bucket
```

## Flags

```
  -a, --access-key string    AWS Access Key
  -b, --bucket-name string   Bucket Name
      --delete               Actually delete objects (default is dry run)
  -e, --endpoint string      Custom S3 Endpoint
  -r, --region string        AWS Region (default "us-east-1")
  -s, --secret-key string    AWS Secret Key
      --yes                  Skip the interactive "yes" confirmation, wait 10s (cancel with Ctrl-C) and delete
```
