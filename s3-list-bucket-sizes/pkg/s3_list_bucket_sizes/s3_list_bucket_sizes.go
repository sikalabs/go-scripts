package s3_list_bucket_sizes

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Options struct {
	AccessKey string
	SecretKey string
	Region    string
	Endpoint  string
}

func S3ListBucketSizes(opts Options) error {
	ctx := context.Background()

	client := s3.New(s3.Options{
		Region:      opts.Region,
		Credentials: credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, ""),
	}, func(o *s3.Options) {
		if opts.Endpoint != "" {
			o.BaseEndpoint = aws.String(opts.Endpoint)
			// S3-compatible services (MinIO, ...) usually require path-style addressing
			o.UsePathStyle = true
		}
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "BUCKET\tOBJECTS\tSIZE")

	var totalObjects, totalSize int64

	paginator := s3.NewListBucketsPaginator(client, &s3.ListBucketsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list buckets: %w", err)
		}
		for _, bucket := range page.Buckets {
			name := aws.ToString(bucket.Name)

			// On AWS, buckets in other regions have to be accessed via their own region
			region := ""
			if opts.Endpoint == "" {
				region = aws.ToString(bucket.BucketRegion)
			}

			objects, size, err := bucketSize(ctx, client, name, region)
			if err != nil {
				fmt.Fprintf(w, "%s\t-\t-\terror: %v\n", name, err)
				continue
			}

			fmt.Fprintf(w, "%s\t%d\t%s\n", name, objects, humanSize(size))
			totalObjects += objects
			totalSize += size
		}
	}

	fmt.Fprintf(w, "TOTAL\t%d\t%s\n", totalObjects, humanSize(totalSize))

	return w.Flush()
}

func bucketSize(ctx context.Context, client *s3.Client, bucket, region string) (int64, int64, error) {
	var objects, size int64

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx, func(o *s3.Options) {
			if region != "" {
				o.Region = region
			}
		})
		if err != nil {
			return 0, 0, err
		}
		for _, obj := range page.Contents {
			objects++
			size += aws.ToInt64(obj.Size)
		}
	}

	return objects, size, nil
}

func humanSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}
