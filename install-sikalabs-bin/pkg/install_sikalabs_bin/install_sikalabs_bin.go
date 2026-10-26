package install_sikalabs_bin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const baseURL = "https://sikalabs.fra1.digitaloceanspaces.com/bin"

func InstallBin(name, dir string) (string, error) {
	url := fmt.Sprintf("%s/%s-linux-amd64", baseURL, name)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download %s: %s", url, resp.Status)
	}

	// Download to a temp file in the target dir and rename it, so an
	// already installed (possibly running) binary is replaced atomically
	tmp, err := os.CreateTemp(dir, "."+name+"-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file in %s: %w", dir, err)
	}
	defer os.Remove(tmp.Name())

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		tmp.Close()
		return "", fmt.Errorf("failed to write %s: %w", tmp.Name(), err)
	}

	err = tmp.Close()
	if err != nil {
		return "", fmt.Errorf("failed to write %s: %w", tmp.Name(), err)
	}

	err = os.Chmod(tmp.Name(), 0755)
	if err != nil {
		return "", fmt.Errorf("failed to chmod %s: %w", tmp.Name(), err)
	}

	dst := filepath.Join(dir, name)
	err = os.Rename(tmp.Name(), dst)
	if err != nil {
		return "", fmt.Errorf("failed to move binary to %s: %w", dst, err)
	}

	return dst, nil
}
