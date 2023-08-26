package installable

import (
	"errors"
	"path"

	"github.com/magefile/mage/sh"
)

const GoBinaryType = "go:binary"

func newGoBinary(entry Entry, opt *GoBinaryOption) *GoBinary {
	return &GoBinary{
		Name:    entry.Name,
		Version: entry.Version,
		Source:  entry.Source,

		option:    *opt,
		versioned: entry.Name + "@" + entry.Version,
	}
}

type GoBinaryOption struct {
	CI string `yaml:"ci"`
}

type GoBinary struct {
	Name    string
	Version string
	Source  string

	option    GoBinaryOption
	versioned string
}

func (g *GoBinary) Install(dst string) (string, error) {
	installed := path.Join(dst, g.versioned)
	if err := checkInstalled(dst, g.Name, g.versioned, g.option.CI); err != nil {
		if errors.Is(err, ErrAlreadyInstalled) {
			return installed, nil
		}
		return installed, err
	}

	env := map[string]string{
		"GOBIN": installed,
	}
	return installed, sh.RunWithV(env, "go", "install", g.Source+"@"+g.Version)
}

func (n *GoBinary) Runtime() Installable {
	return nil
}
