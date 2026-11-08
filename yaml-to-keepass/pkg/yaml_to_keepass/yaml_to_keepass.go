package yaml_to_keepass

import (
	"fmt"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	MasterPassword string  `yaml:"master_password"`
	Entries        []Entry `yaml:"entries"`
	Groups         []Group `yaml:"groups"`
}

type Group struct {
	Name    string  `yaml:"name"`
	Entries []Entry `yaml:"entries"`
	Groups  []Group `yaml:"groups"`
}

type Entry struct {
	Title    string `yaml:"title"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
	Notes    string `yaml:"notes"`
}

func readConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}
	if config.MasterPassword == "" {
		config.MasterPassword = os.Getenv("KEEPASS_MASTER_PASSWORD")
	}
	if config.MasterPassword == "" {
		return nil, fmt.Errorf("master password is required (master_password in YAML or KEEPASS_MASTER_PASSWORD env var)")
	}
	return &config, nil
}

func YamlToKeepass(inputPath, outputPath string) error {
	config, err := readConfig(inputPath)
	if err != nil {
		return err
	}

	root := gokeepasslib.NewGroup()
	root.Name = "Root"
	root.Entries = buildEntries(config.Entries)
	root.Groups = buildGroups(config.Groups)

	db := gokeepasslib.NewDatabase(gokeepasslib.WithDatabaseKDBXVersion4())
	db.Credentials = gokeepasslib.NewPasswordCredentials(config.MasterPassword)
	db.Content.Root.Groups = []gokeepasslib.Group{root}

	if err := db.LockProtectedEntries(); err != nil {
		return err
	}

	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gokeepasslib.NewEncoder(file).Encode(db); err != nil {
		return err
	}
	return file.Close()
}

func buildGroups(groups []Group) []gokeepasslib.Group {
	out := make([]gokeepasslib.Group, 0, len(groups))
	for _, g := range groups {
		group := gokeepasslib.NewGroup()
		group.Name = g.Name
		group.Entries = buildEntries(g.Entries)
		group.Groups = buildGroups(g.Groups)
		out = append(out, group)
	}
	return out
}

func buildEntries(entries []Entry) []gokeepasslib.Entry {
	out := make([]gokeepasslib.Entry, 0, len(entries))
	for _, e := range entries {
		entry := gokeepasslib.NewEntry()
		entry.Values = append(entry.Values,
			mkValue("Title", e.Title),
			mkValue("UserName", e.Username),
			mkProtectedValue("Password", e.Password),
			mkValue("URL", e.URL),
			mkValue("Notes", e.Notes),
		)
		out = append(out, entry)
	}
	return out
}

func mkValue(key, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}}
}

func mkProtectedValue(key, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value, Protected: w.NewBoolWrapper(true)},
	}
}
