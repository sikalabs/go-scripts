package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/sikalabs/go-scripts/s3-upload-file/pkg/s3_upload_file"
	"github.com/spf13/cobra"
)

var FlagFile string
var FlagBucket string
var FlagKey string
var FlagRegion string
var FlagEndpoint string
var FlagPathStyle bool
var FlagACL string
var FlagPublic bool
var FlagAccessKey string
var FlagSecretKey string

var Cmd = &cobra.Command{
	Use:   "s3-upload-file",
	Short: "Upload a file to an S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		filePath := FlagFile

		key := FlagKey
		if key == "" {
			key = filepath.Base(filePath)
		}

		acl := FlagACL
		if FlagPublic {
			acl = "public-read"
		}

		err := s3_upload_file.UploadFile(s3_upload_file.Options{
			FilePath:  filePath,
			Bucket:    FlagBucket,
			Key:       key,
			Region:    FlagRegion,
			Endpoint:  FlagEndpoint,
			PathStyle: FlagPathStyle,
			ACL:       acl,
			AccessKey: FlagAccessKey,
			SecretKey: FlagSecretKey,
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Uploaded %s to s3://%s/%s\n", filePath, FlagBucket, key)
	},
}

func init() {
	Cmd.Flags().StringVarP(
		&FlagFile,
		"file",
		"f",
		"",
		"Path of the file to upload",
	)
	Cmd.MarkFlagRequired("file")

	Cmd.Flags().StringVarP(
		&FlagBucket,
		"bucket",
		"b",
		"",
		"S3 bucket name",
	)
	Cmd.MarkFlagRequired("bucket")

	Cmd.Flags().StringVarP(
		&FlagKey,
		"key",
		"k",
		"",
		"S3 object key (defaults to the file name)",
	)

	Cmd.Flags().StringVarP(
		&FlagRegion,
		"region",
		"r",
		"us-east-1",
		"AWS region",
	)

	Cmd.Flags().StringVar(
		&FlagEndpoint,
		"endpoint",
		"",
		"Custom S3 endpoint URL (for S3-compatible services)",
	)

	Cmd.Flags().BoolVar(
		&FlagPathStyle,
		"path-style",
		false,
		"Use path-style addressing (required by some S3-compatible services)",
	)

	Cmd.Flags().StringVar(
		&FlagACL,
		"acl",
		"",
		"Canned ACL to apply to the uploaded object (e.g. private, public-read)",
	)

	Cmd.Flags().BoolVar(
		&FlagPublic,
		"public",
		false,
		"Shortcut for --acl public-read",
	)

	Cmd.Flags().StringVar(
		&FlagAccessKey,
		"access-key",
		"",
		"AWS/S3 access key (falls back to standard AWS credential chain if unset)",
	)

	Cmd.Flags().StringVar(
		&FlagSecretKey,
		"secret-key",
		"",
		"AWS/S3 secret key (falls back to standard AWS credential chain if unset)",
	)
}
