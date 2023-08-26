package installable

import (
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type Installables map[string]Installable

func LoadInstallables(data []byte) (Installables, error) {
	loaded := new(Entries)
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		return nil, err
	}

	installables := make(Installables, len(loaded.Data))

	for _, entry := range loaded.Data {
		switch entry.Type {
		case HTTPArchiveType:
			opt, err := typedOption[HTTPArchiveOption](entry)
			if err != nil {
				return nil, err
			}
			installables[entry.Name] = newHTTPArchive(entry, opt)
		case NPMBinaryType:
			opt, err := typedOption[NpmBinaryOption](entry)
			if err != nil {
				return nil, err
			}
			if opt.Runtime != "" {
				// By default, we surpress the error. We will use system installed runtime when we run this binary.
				ri, _ := loaded.Resolve(opt.Runtime)
				opt.RuntimeInstaller = ri
			}
			installables[entry.Name] = newNPMBinary(entry, opt)
		case GoBinaryType:
			opt, err := typedOption[GoBinaryOption](entry)
			if err != nil {
				return nil, err
			}
			installables[entry.Name] = newGoBinary(entry, opt)
		}
	}
	return installables, nil
}

type Installable interface {
	Install(string) (string, error)
	Runtime() Installable
}

type InstallabeInfo struct {
	Key        string
	Binary     string
	Installers []Installable
}

type option interface {
	HTTPArchiveOption | NpmBinaryOption | GoBinaryOption
}

type Entry struct {
	Name    string      `yaml:"name"`
	Version string      `yaml:"version"`
	Source  string      `yaml:"source"`
	Type    string      `yaml:"type"`
	Option  interface{} `yaml:"option"`
}

type Entries struct {
	Data []Entry `yaml:"tools"`
}

func (e *Entries) Resolve(name string) (Installable, error) {
	for _, entry := range e.Data {
		if entry.Name != name {
			continue
		}
		switch entry.Type {
		case "http:archive":
			opt, err := typedOption[HTTPArchiveOption](entry)
			if err != nil {
				return nil, err
			}
			return newHTTPArchive(entry, opt), nil
		case "npm:binary":
			opt, err := typedOption[NpmBinaryOption](entry)
			if err != nil {
				return nil, err
			}
			return newNPMBinary(entry, opt), nil
		case "go:binary":
			opt, err := typedOption[GoBinaryOption](entry)
			if err != nil {
				return nil, err
			}
			return newGoBinary(entry, opt), nil
		}
	}
	return nil, ErrNotFound
}

func fromUntypedOption[T option](i interface{}, typed *T) error {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, typed)
}

func typedOption[T option](e Entry) (*T, error) {
	opt := new(T)
	if err := fromUntypedOption[T](e.Option, opt); err != nil {
		return nil, err
	}
	return opt, nil
}

func GetInstalledBase(installed string) string {
	parts := strings.Split(installed, "/")
	var base []string
	for _, part := range parts {
		base = append(base, part)
		if strings.Contains(part, "@v") {
			return path.Join(base...)
		}
	}
	return installed
}

// TODO(dio): Make sure when writing to dst, we always have "v"-prefix for versions.
func checkInstalled(dir, prefix, current, ci string) error {
	// For example in GHA, we often use @actions/setup-node so we can skip installing.
	if ci == "skip" && os.Getenv("CI") == "true" {
		return ErrAlreadyInstalled
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name() == current { // TODO(dio): Check content.
			return ErrAlreadyInstalled
		}

		if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix+"@") {
			// TODO(dio): Remove this when we need/allow multiple versions. Note that we also need
			// to have a querier (for sorting paths with the right order, or pointing to the right binary).
			if err := os.RemoveAll(path.Join(dir, entry.Name())); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func infer(m map[string]string, ref, fallback string) string {
	for k, v := range m {
		if k == ref {
			return v
		}
	}
	return fallback
}
