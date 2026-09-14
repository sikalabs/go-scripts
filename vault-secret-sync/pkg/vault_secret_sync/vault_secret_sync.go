package vault_secret_sync

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/vault/api"
)

// mount describes the secrets engine that owns a given path, so reads and
// writes can be translated to the right KV v1/v2 endpoints.
type mount struct {
	path    string // mount path, e.g. "secret/"
	version string // "1" or "2"
	sub     string // path relative to the mount
}

func newClient(addr, token string) *api.Client {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Vault client for %s: %v", addr, err)
	}
	client.SetToken(token)
	return client
}

func findMount(client *api.Client, path string) (*mount, error) {
	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return nil, fmt.Errorf("failed to list mounts: %w", err)
	}

	path = strings.TrimPrefix(path, "/")

	var best string
	for mountPath, m := range mounts {
		if m.Type != "kv" && m.Type != "generic" {
			continue
		}
		if path == strings.TrimSuffix(mountPath, "/") || strings.HasPrefix(path, mountPath) {
			if len(mountPath) > len(best) {
				best = mountPath
			}
		}
	}
	if best == "" {
		return nil, fmt.Errorf("no KV mount found for path %q", path)
	}

	version := "1"
	if v, ok := mounts[best].Options["version"]; ok && v == "2" {
		version = "2"
	}

	return &mount{
		path:    best,
		version: version,
		sub:     strings.TrimPrefix(path, best),
	}, nil
}

func listPath(client *api.Client, m *mount, sub string) ([]string, error) {
	p := m.path + sub
	if m.version == "2" {
		p = m.path + "metadata/" + sub
	}

	secret, err := client.Logical().List(p)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	rawKeys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, nil
	}

	keys := make([]string, 0, len(rawKeys))
	for _, k := range rawKeys {
		keys = append(keys, k.(string))
	}
	return keys, nil
}

func readSecret(client *api.Client, m *mount, sub string) (map[string]interface{}, error) {
	readPath := m.path + sub
	if m.version == "2" {
		readPath = m.path + "data/" + sub
	}

	secret, err := client.Logical().Read(readPath)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	if m.version == "2" {
		data, ok := secret.Data["data"].(map[string]interface{})
		if !ok || data == nil {
			return nil, nil
		}
		return data, nil
	}

	return secret.Data, nil
}

func writeSecret(client *api.Client, m *mount, sub string, data map[string]interface{}) error {
	writePath := m.path + sub
	body := data
	if m.version == "2" {
		writePath = m.path + "data/" + sub
		body = map[string]interface{}{"data": data}
	}

	_, err := client.Logical().Write(writePath, body)
	return err
}

func syncPath(sourceClient *api.Client, sourceMount *mount, targetClient *api.Client, targetMount *mount, sub string) {
	keys, err := listPath(sourceClient, sourceMount, sub)
	if err != nil {
		log.Fatalf("Failed to list %q: %v", sub, err)
	}

	if len(keys) == 0 {
		data, err := readSecret(sourceClient, sourceMount, sub)
		if err != nil {
			log.Fatalf("Failed to read secret %q: %v", sub, err)
		}
		if data == nil {
			return
		}
		fmt.Printf("Syncing secret: %s\n", sub)
		if err := writeSecret(targetClient, targetMount, sub, data); err != nil {
			log.Fatalf("Failed to write secret %q: %v", sub, err)
		}
		return
	}

	for _, key := range keys {
		childSub := sub + key
		if strings.HasSuffix(key, "/") {
			syncPath(sourceClient, sourceMount, targetClient, targetMount, childSub)
			continue
		}

		data, err := readSecret(sourceClient, sourceMount, childSub)
		if err != nil {
			log.Fatalf("Failed to read secret %q: %v", childSub, err)
		}
		if data == nil {
			continue
		}
		fmt.Printf("Syncing secret: %s\n", childSub)
		if err := writeSecret(targetClient, targetMount, childSub, data); err != nil {
			log.Fatalf("Failed to write secret %q: %v", childSub, err)
		}
	}
}

func VaultSecretSync(sourceAddr, sourceToken, targetAddr, targetToken, path string) {
	sourceClient := newClient(sourceAddr, sourceToken)
	targetClient := newClient(targetAddr, targetToken)

	sourceMount, err := findMount(sourceClient, path)
	if err != nil {
		log.Fatalf("Failed to resolve source mount: %v", err)
	}
	targetMount, err := findMount(targetClient, path)
	if err != nil {
		log.Fatalf("Failed to resolve target mount: %v", err)
	}

	sub := strings.TrimSuffix(sourceMount.sub, "/")
	if sub != "" {
		sub += "/"
	}

	syncPath(sourceClient, sourceMount, targetClient, targetMount, sub)

	fmt.Println("Done!")
}
