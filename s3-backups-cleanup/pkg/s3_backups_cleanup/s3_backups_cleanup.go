package s3_backups_cleanup

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const KeepDailyDays = 30

// How long to wait before deleting with --yes, gives time to cancel with Ctrl-C
const YesDelay = 10 * time.Second

// Matches the date in backup names like sikalabs-com.2026-10-22_18-05-09_mqk.sql.gz
var backupDateRegexp = regexp.MustCompile(`\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}`)

const backupDateLayout = "2006-01-02_15-04-05"

type Options struct {
	AccessKey  string
	SecretKey  string
	Region     string
	Endpoint   string
	BucketName string
	Delete     bool
	Yes        bool
}

type Object struct {
	Key string
	// Name is the part of the key before the date (e.g. "sikalabs-com."),
	// retention is applied to each name separately
	Name string
	Date time.Time
}

func S3BackupsCleanup(opts Options) error {
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

	objects, skipped, err := listObjects(ctx, client, opts.BucketName)
	if err != nil {
		return err
	}

	for _, key := range skipped {
		fmt.Printf("SKIP   %s (no date in name)\n", key)
	}

	keep := objectsToKeep(objects, time.Now())

	var toDelete []Object
	for _, obj := range objects {
		if keep[obj.Key] {
			fmt.Printf("KEEP   %s  %s\n", obj.Date.Format(time.RFC3339), obj.Key)
			continue
		}
		fmt.Printf("DELETE %s  %s\n", obj.Date.Format(time.RFC3339), obj.Key)
		toDelete = append(toDelete, obj)
	}

	fmt.Printf("\nKeep %d, delete %d objects\n", len(keep), len(toDelete))

	if !opts.Delete {
		fmt.Println("Dry run, use --delete to delete")
		return nil
	}

	if len(toDelete) == 0 {
		return nil
	}

	if opts.Yes {
		if !waitForCancel() {
			fmt.Println("Delete cancelled.")
			return nil
		}
	} else if !confirm() {
		fmt.Println("Delete cancelled.")
		return nil
	}

	for _, obj := range toDelete {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(opts.BucketName),
			Key:    aws.String(obj.Key),
		})
		if err != nil {
			return fmt.Errorf("failed to delete s3://%s/%s: %w", opts.BucketName, obj.Key, err)
		}
		fmt.Printf("Deleted %s\n", obj.Key)
	}

	fmt.Printf("\nDeleted %d objects\n", len(toDelete))

	return nil
}

// confirm asks the user to type "yes" (like terraform apply)
func confirm() bool {
	fmt.Println()
	fmt.Println("Do you want to delete these objects?")
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

// listObjects returns objects with a date in the name and keys of objects without it
func listObjects(ctx context.Context, client *s3.Client, bucket string) ([]Object, []string, error) {
	var objects []Object
	var skipped []string

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list objects in s3://%s: %w", bucket, err)
		}
		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			name, date, ok := parseBackupKey(key)
			if !ok {
				skipped = append(skipped, key)
				continue
			}
			objects = append(objects, Object{
				Key:  key,
				Name: name,
				Date: date,
			})
		}
	}

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].Date.Before(objects[j].Date)
	})

	return objects, skipped, nil
}

func parseBackupKey(key string) (string, time.Time, bool) {
	loc := backupDateRegexp.FindStringIndex(key)
	if loc == nil {
		return "", time.Time{}, false
	}
	date, err := time.Parse(backupDateLayout, key[loc[0]:loc[1]])
	if err != nil {
		return "", time.Time{}, false
	}
	return key[:loc[0]], date, true
}

// objectsToKeep expects objects sorted from the oldest to the newest.
// For every name it keeps the first (oldest) object of every month forever
// and the last (newest) object of every day in the last KeepDailyDays days.
func objectsToKeep(objects []Object, now time.Time) map[string]bool {
	keep := map[string]bool{}
	dailyCutoff := now.UTC().AddDate(0, 0, -KeepDailyDays)

	monthly := map[string]string{}
	daily := map[string]string{}

	for _, obj := range objects {
		month := obj.Name + "|" + obj.Date.Format("2006-01")
		if _, ok := monthly[month]; !ok {
			monthly[month] = obj.Key
		}

		if obj.Date.After(dailyCutoff) {
			daily[obj.Name+"|"+obj.Date.Format("2006-01-02")] = obj.Key
		}
	}

	for _, key := range monthly {
		keep[key] = true
	}
	for _, key := range daily {
		keep[key] = true
	}

	return keep
}
