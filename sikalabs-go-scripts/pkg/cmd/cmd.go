package cmd

import (
	digitalocean_cleanup_empty_loadbalancers "github.com/sikalabs/go-scripts/digitalocean-cleanup-empty-loadbalancers/pkg/cmd"
	digitalocean_cleanup_unused_volumes "github.com/sikalabs/go-scripts/digitalocean-cleanup-unused-volumes/pkg/cmd"
	example_hello_world "github.com/sikalabs/go-scripts/example-hello-world/pkg/cmd"
	get_secret_from_keepass "github.com/sikalabs/go-scripts/get-secret-from-keepass/pkg/cmd"
	install_sikalabs_bin "github.com/sikalabs/go-scripts/install-sikalabs-bin/pkg/cmd"
	keepass_to_yaml "github.com/sikalabs/go-scripts/keepass-to-yaml/pkg/cmd"
	kube_port_forward_multi "github.com/sikalabs/go-scripts/kube-port-forward-multi/pkg/cmd"
	mirror_container_registry "github.com/sikalabs/go-scripts/mirror-container-registry/pkg/cmd"
	s3_backups_cleanup "github.com/sikalabs/go-scripts/s3-backups-cleanup/pkg/cmd"
	s3_bucket_cleanup "github.com/sikalabs/go-scripts/s3-bucket-cleanup/pkg/cmd"
	s3_list_bucket_sizes "github.com/sikalabs/go-scripts/s3-list-bucket-sizes/pkg/cmd"
	s3_upload_file "github.com/sikalabs/go-scripts/s3-upload-file/pkg/cmd"
	telegram_send_message "github.com/sikalabs/go-scripts/telegram-send-message/pkg/cmd"
	upload_sikalabs_bin "github.com/sikalabs/go-scripts/upload-sikalabs-bin/pkg/cmd"
	vault_secret_sync "github.com/sikalabs/go-scripts/vault-secret-sync/pkg/cmd"
	yaml_to_keepass "github.com/sikalabs/go-scripts/yaml-to-keepass/pkg/cmd"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "sikalabs-go-scripts",
	Short: "All SikaLabs Go scripts in one binary",
}

func init() {
	Cmd.AddCommand(
		digitalocean_cleanup_empty_loadbalancers.Cmd,
		digitalocean_cleanup_unused_volumes.Cmd,
		example_hello_world.Cmd,
		get_secret_from_keepass.Cmd,
		install_sikalabs_bin.Cmd,
		keepass_to_yaml.Cmd,
		kube_port_forward_multi.Cmd,
		mirror_container_registry.Cmd,
		s3_backups_cleanup.Cmd,
		s3_bucket_cleanup.Cmd,
		s3_list_bucket_sizes.Cmd,
		s3_upload_file.Cmd,
		telegram_send_message.Cmd,
		upload_sikalabs_bin.Cmd,
		vault_secret_sync.Cmd,
		yaml_to_keepass.Cmd,
	)
}
