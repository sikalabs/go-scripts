package mirror_git_repo

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
)

func MirrorGitRepo(sourceUrl, targetUrl, sourceToken, targetToken string) {
	if sourceUrl == "" {
		log.Fatal("--source-url is required")
	}
	if targetUrl == "" {
		log.Fatal("--target-url is required")
	}

	sourceAuthUrl, err := withAuth(sourceUrl, sourceToken)
	if err != nil {
		log.Fatalf("Failed to parse --source-url: %v", err)
	}
	targetAuthUrl, err := withAuth(targetUrl, targetToken)
	if err != nil {
		log.Fatalf("Failed to parse --target-url: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "mirror-git-repo-*")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Printf("Cloning %s\n", sourceUrl)
	if err := runGit("", "clone", "--mirror", sourceAuthUrl, tmpDir); err != nil {
		log.Fatalf("Failed to clone source repo: %v", err)
	}

	fmt.Printf("Pushing to %s\n", targetUrl)
	if err := runGit(tmpDir, "push", "--mirror", targetAuthUrl); err != nil {
		log.Fatalf("Failed to push to target repo: %v", err)
	}

	fmt.Println("Done!")
}

// withAuth returns rawUrl with token injected as the URL user info, so git
// can authenticate over HTTPS without prompting or needing a credential
// helper. Non-HTTP(S) URLs (e.g. git@host:path.git) are returned unchanged.
func withAuth(rawUrl, token string) (string, error) {
	if token == "" {
		return rawUrl, nil
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return rawUrl, nil
	}

	u.User = url.User(token)
	return u.String(), nil
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	return cmd.Run()
}
