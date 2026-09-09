package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
)

// registryKeychain resolves credentials per-registry host, since the source
// and target registries use different deploy tokens and crane.Copy applies
// a single Option set to both the pull and the push.
type registryKeychain struct {
	sourceRegistry string
	sourceAuth     authn.Authenticator
	targetRegistry string
	targetAuth     authn.Authenticator
}

func (k registryKeychain) Resolve(target authn.Resource) (authn.Authenticator, error) {
	switch target.RegistryStr() {
	case k.sourceRegistry:
		return k.sourceAuth, nil
	case k.targetRegistry:
		return k.targetAuth, nil
	default:
		return authn.Anonymous, nil
	}
}

func main() {
	sourceRegistry := flag.String("source", "", "Source registry (e.g. registry-source.example.com)")
	targetRegistry := flag.String("target", "", "Target registry (e.g. registry-target.example.com)")
	sourceGroup := flag.String("source-group", "", "Source group to filter (e.g. source-group)")
	targetGroup := flag.String("target-group", "", "Target group, replaces source group (e.g. target-group)")
	sourceToken := flag.String("source-token", "", "Source registry deploy token (optional)")
	targetToken := flag.String("target-token", "", "Target registry deploy token (optional)")
	flag.Parse()

	if *sourceRegistry == "" {
		log.Fatal("--source is required")
	}
	if *targetRegistry == "" {
		log.Fatal("--target is required")
	}
	if *sourceGroup == "" {
		log.Fatal("--source-group is required")
	}
	if *targetGroup == "" {
		log.Fatal("--target-group is required")
	}

	sourceAuth := authn.Anonymous
	if *sourceToken != "" {
		sourceAuth = &authn.Basic{Username: "token", Password: *sourceToken}
	}
	targetAuth := authn.Anonymous
	if *targetToken != "" {
		targetAuth = &authn.Basic{Username: "token", Password: *targetToken}
	}

	kc := registryKeychain{
		sourceRegistry: *sourceRegistry,
		sourceAuth:     sourceAuth,
		targetRegistry: *targetRegistry,
		targetAuth:     targetAuth,
	}
	authOpt := crane.WithAuthFromKeychain(kc)

	// Get all repos from source registry
	fmt.Printf("Fetching catalog from %s\n", *sourceRegistry)
	repos, err := crane.Catalog(*sourceRegistry, authOpt)
	if err != nil {
		log.Fatalf("Failed to get catalog: %v", err)
	}

	// Sync matching repos
	for _, repo := range repos {
		if strings.HasPrefix(repo, *sourceGroup) {
			syncRepo(*sourceRegistry, *targetRegistry, repo, *sourceGroup, *targetGroup, authOpt)
		}
	}

	fmt.Println("Done!")
}

func syncRepo(sourceRegistry, targetRegistry, repo, sourceGroup, targetGroup string, authOpt crane.Option) {
	fmt.Printf("Syncing repo: %s\n", repo)

	// List tags
	tags, err := crane.ListTags(fmt.Sprintf("%s/%s", sourceRegistry, repo), authOpt)
	if err != nil {
		log.Printf("Warning: failed to list tags for %s: %v", repo, err)
		return
	}

	// Replace source group with target group in path
	targetRepo := strings.Replace(repo, sourceGroup, targetGroup, 1)

	for _, tag := range tags {
		src := fmt.Sprintf("%s/%s:%s", sourceRegistry, repo, tag)
		dst := fmt.Sprintf("%s/%s:%s", targetRegistry, targetRepo, tag)

		fmt.Printf("  Copying %s -> %s\n", src, dst)

		if err := crane.Copy(src, dst, authOpt); err != nil {
			log.Printf("Error copying %s: %v", src, err)
		}
	}
}
