package get_secret_from_keepass

import (
	"fmt"
	"os"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
)

func GetSecretFromKeepass(dbPath, secretPath, password string) (string, error) {
	file, err := os.Open(dbPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)
	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}
	if err := db.UnlockProtectedEntries(); err != nil {
		return "", err
	}

	entry, err := findEntry(db.Content.Root.Groups, secretPath)
	if err != nil {
		return "", err
	}
	return entry.GetPassword(), nil
}

// findEntry finds an entry by path like /Root/group1/bar_entry,
// where all parts except the last are group names and the last
// part is the entry title.
func findEntry(groups []gokeepasslib.Group, secretPath string) (*gokeepasslib.Entry, error) {
	parts := strings.Split(strings.Trim(secretPath, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid path %q, expected /<group>/.../<entry>", secretPath)
	}
	groupNames, entryTitle := parts[:len(parts)-1], parts[len(parts)-1]

	var group *gokeepasslib.Group
	for _, name := range groupNames {
		group = nil
		for i := range groups {
			if groups[i].Name == name {
				group = &groups[i]
				break
			}
		}
		if group == nil {
			return nil, fmt.Errorf("group %q not found in path %q", name, secretPath)
		}
		groups = group.Groups
	}

	for i := range group.Entries {
		if group.Entries[i].GetTitle() == entryTitle {
			return &group.Entries[i], nil
		}
	}
	return nil, fmt.Errorf("entry %q not found in path %q", entryTitle, secretPath)
}
