package s3_bucket_cleanup

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// How long to wait before deleting with --yes, gives time to cancel with Ctrl-C
const YesDelay = 10 * time.Second

type Options struct {
	AccessKey  string
	SecretKey  string
	Region     string
	Endpoint   string
	BucketName string
	Delete     bool
	Yes        bool
}

func S3BucketCleanup(opts Options) error {
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

	keys, err := listKeys(ctx, client, opts.BucketName)
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Printf("DELETE %s\n", key)
	}

	fmt.Printf("\nDelete %d objects from s3://%s\n", len(keys), opts.BucketName)

	if !opts.Delete {
		fmt.Println("Dry run, use --delete to delete")
		return nil
	}

	if len(keys) == 0 {
		return nil
	}

	if opts.Yes {
		if !waitForCancel() {
			fmt.Println("Delete cancelled.")
			return nil
		}
	} else if !confirm(opts.BucketName) {
		fmt.Println("Delete cancelled.")
		return nil
	}

	for _, key := range keys {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(opts.BucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("failed to delete s3://%s/%s: %w", opts.BucketName, key, err)
		}
		fmt.Printf("Deleted %s\n", key)
	}

	fmt.Printf("\nDeleted %d objects\n", len(keys))

	return nil
}

func listKeys(ctx context.Context, client *s3.Client, bucket string) ([]string, error) {
	var keys []string

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects in s3://%s: %w", bucket, err)
		}
		for _, obj := range page.Contents {
			keys = append(keys, aws.ToString(obj.Key))
		}
	}

	return keys, nil
}

// confirm asks the user to type "yes" (like terraform apply)
func confirm(bucket string) bool {
	fmt.Println()
	fmt.Printf("Do you want to delete ALL objects in s3://%s?\n", bucket)
	fmt.Println("  Only 'yes' will be accepted to approve.")
	fmt.Println()
	fmt.Print("  Enter a value: ")

	answer, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println()

	return strings.TrimSpace(answer) == "yes"
}

// waitForCancel waits YesDelay and returns false if the user pressed Ctrl-C
func waitForCancel() bool {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	defer signal.Stop(sigs)

	fmt.Printf("\nDeleting in %s, press Ctrl-C to cancel...\n", YesDelay)

	select {
	case <-sigs:
		fmt.Println()
		return false
	case <-time.After(YesDelay):
		return true
	}
}
