package keepass_to_yaml

import (
	"bytes"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
	"go.yaml.in/yaml/v3"
)

const FormatVersion = 1

type Metadata struct {
	Version int `yaml:"version"`
}

type Config struct {
	Metadata       Metadata `yaml:"metadata"`
	MasterPassword string   `yaml:"master_password"`
	Entries        []Entry  `yaml:"entries,omitempty"`
	Groups         []Group  `yaml:"groups,omitempty"`
}

type Group struct {
	Name    string  `yaml:"name"`
	Entries []Entry `yaml:"entries,omitempty"`
	Groups  []Group `yaml:"groups,omitempty"`
}

type Entry struct {
	Title    string `yaml:"title,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	URL      string `yaml:"url,omitempty"`
	Notes    string `yaml:"notes,omitempty"`
}

func readDatabase(path, password string) (*gokeepasslib.Database, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)
	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		return nil, err
	}
	if err := db.UnlockProtectedEntries(); err != nil {
		return nil, err
	}
	return db, nil
}

func KeepassToYaml(inputPath, outputPath, password string) error {
	db, err := readDatabase(inputPath, password)
	if err != nil {
		return err
	}

	var recycleBin *gokeepasslib.UUID
	if db.Content.Meta != nil && db.Content.Meta.RecycleBinEnabled.Bool {
		recycleBin = &db.Content.Meta.RecycleBinUUID
	}

	config := Config{
		Metadata:       Metadata{Version: FormatVersion},
		MasterPassword: password,
	}
	roots := db.Content.Root.Groups
	if len(roots) == 1 {
		// The single root group (e.g. "Root") maps to the top level of the YAML.
		config.Entries = buildEntries(roots[0].Entries)
		config.Groups = buildGroups(roots[0].Groups, recycleBin)
	} else {
		config.Groups = buildGroups(roots, recycleBin)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&config); err != nil {
		return err
	}
	if err := encoder.Close(); err != nil {
		return err
	}
	return os.WriteFile(outputPath, buf.Bytes(), 0600)
}

func buildGroups(groups []gokeepasslib.Group, recycleBin *gokeepasslib.UUID) []Group {
	out := make([]Group, 0, len(groups))
	for _, g := range groups {
		if recycleBin != nil && g.UUID.Compare(*recycleBin) {
			continue
		}
		out = append(out, Group{
			Name:    g.Name,
			Entries: buildEntries(g.Entries),
			Groups:  buildGroups(g.Groups, recycleBin),
		})
	}
	return out
}

func buildEntries(entries []gokeepasslib.Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, Entry{
			Title:    e.GetTitle(),
			Username: e.GetContent("UserName"),
			Password: e.GetPassword(),
			URL:      e.GetContent("URL"),
			Notes:    e.GetContent("Notes"),
		})
	}
	return out
}
