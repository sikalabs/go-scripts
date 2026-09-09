# s3-upload-file

Upload a file to an S3 bucket.

## Install (master)

```
go install github.com/sikalabs/go-scripts/s3-upload-file@master
```

## Example Usage

```bash
s3-upload-file --file ./file.txt --bucket my-bucket
```

or with a custom key

```bash
s3-upload-file --file ./file.txt --bucket my-bucket --key path/in/bucket/file.txt
```

### S3-compatible services (e.g. MinIO)

```bash
s3-upload-file \
  --file ./file.txt \
  --bucket my-bucket \
  --endpoint https://minio.example.com \
  --path-style
```

## Credentials

AWS credentials and region are resolved using the standard AWS SDK
credential chain (environment variables, shared config/credentials
files, IAM role, etc.). Use `--region` to override the region.
