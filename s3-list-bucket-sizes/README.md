# s3-list-bucket-sizes

List all S3 buckets with their number of objects and total size.

The size is computed by listing all objects in each bucket, so it can take
a while for buckets with many objects. Old versions of objects in versioned
buckets are not counted.

## Install (master)

```
go install github.com/sikalabs/go-scripts/s3-list-bucket-sizes@master
```

## Example Usage

```bash
s3-list-bucket-sizes \
  --access-key xxx \
  --secret-key xxx
```

```
BUCKET      OBJECTS  SIZE
my-backups  284      1.2 GiB
my-bucket   12       34.5 KiB
TOTAL       296      1.2 GiB
```

### S3-compatible services (e.g. MinIO, DigitalOcean Spaces)

Path-style addressing is used automatically with a custom endpoint.

```bash
s3-list-bucket-sizes \
  --access-key xxx \
  --secret-key xxx \
  --endpoint https://fra1.digitaloceanspaces.com \
  --region fra1
```

## Flags

```
  -a, --access-key string   AWS Access Key
  -e, --endpoint string     Custom S3 Endpoint
  -r, --region string       AWS Region (default "us-east-1")
  -s, --secret-key string   AWS Secret Key
```
