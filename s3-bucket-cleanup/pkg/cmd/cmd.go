package cmd

import (
	"fmt"
	"os"

	"github.com/sikalabs/go-scripts/s3-bucket-cleanup/pkg/s3_bucket_cleanup"
	"github.com/spf13/cobra"
)

var FlagAccessKey string
var FlagSecretKey string
var FlagRegion string
var FlagEndpoint string
var FlagBucketName string
var FlagDelete bool
var FlagYes bool

var Cmd = &cobra.Command{
	Use:   "s3-bucket-cleanup",
	Short: "Delete every object in an S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		err := s3_bucket_cleanup.S3BucketCleanup(s3_bucket_cleanup.Options{
			AccessKey:  FlagAccessKey,
			SecretKey:  FlagSecretKey,
			Region:     FlagRegion,
			Endpoint:   FlagEndpoint,
			BucketName: FlagBucketName,
			Delete:     FlagDelete,
			Yes:        FlagYes,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	Cmd.Flags().StringVarP(
		&FlagAccessKey,
		"access-key",
		"a",
		"",
		"AWS Access Key",
	)
	Cmd.MarkFlagRequired("access-key")

	Cmd.Flags().StringVarP(
		&FlagSecretKey,
		"secret-key",
		"s",
		"",
		"AWS Secret Key",
	)
	Cmd.MarkFlagRequired("secret-key")

	Cmd.Flags().StringVarP(
		&FlagRegion,
		"region",
		"r",
		"us-east-1",
		"AWS Region",
	)

	Cmd.Flags().StringVarP(
		&FlagEndpoint,
		"endpoint",
		"e",
		"",
		"Custom S3 Endpoint",
	)

	Cmd.Flags().StringVarP(
		&FlagBucketName,
		"bucket-name",
		"b",
		"",
		"Bucket Name",
	)
	Cmd.MarkFlagRequired("bucket-name")

	Cmd.Flags().BoolVar(
		&FlagDelete,
		"delete",
		false,
		"Actually delete objects (default is dry run)",
	)

	Cmd.Flags().BoolVar(
		&FlagYes,
		"yes",
		false,
		"Skip the interactive \"yes\" confirmation, wait 10s (cancel with Ctrl-C) and delete",
	)
}
