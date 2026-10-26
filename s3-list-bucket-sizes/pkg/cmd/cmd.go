package cmd

import (
	"fmt"
	"os"

	"github.com/sikalabs/go-scripts/s3-list-bucket-sizes/pkg/s3_list_bucket_sizes"
	"github.com/spf13/cobra"
)

var FlagAccessKey string
var FlagSecretKey string
var FlagRegion string
var FlagEndpoint string

var Cmd = &cobra.Command{
	Use:   "s3-list-bucket-sizes",
	Short: "List all S3 buckets with their size and number of objects",
	Run: func(cmd *cobra.Command, args []string) {
		err := s3_list_bucket_sizes.S3ListBucketSizes(s3_list_bucket_sizes.Options{
			AccessKey: FlagAccessKey,
			SecretKey: FlagSecretKey,
			Region:    FlagRegion,
			Endpoint:  FlagEndpoint,
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
}
