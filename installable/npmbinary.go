package installable

import (
	"errors"
	"path"

	"github.com/magefile/mage/sh"
)

const NPMBinaryType = "npm:binary"

func newNPMBinary(entry Entry, opt *NpmBinaryOption) *NpmBinary {
	bin := &NpmBinary{
		Name:    entry.Name,
		Version: entry.Version,
		Source:  entry.Source,

		versioned: entry.Name + "@" + entry.Version,
		option:    *opt,
	}
	if opt.Runtime != "" {
		bin.runtime = opt.RuntimeInstaller
	}
	return bin
}

type NpmBinaryOption struct {
	Runtime          string      `yaml:"runtime"`
	RuntimeInstaller Installable `yaml:"-"`
	CI               string      `yaml:"ci"`
}

type NpmBinary struct {
	Name    string
	Version string
	Source  string

	option NpmBinaryOption

	runtime   Installable
	versioned string
}

func (n *NpmBinary) Runtime() Installable {
	return n.runtime
}

func (n *NpmBinary) Install(dst string) (string, error) {
	installed := path.Join(dst, n.versioned, "node_modules", ".bin")

	if err := checkInstalled(dst, n.Name, n.versioned, n.option.CI); err != nil {
		if errors.Is(err, ErrAlreadyInstalled) {
			return installed, nil
		}
		return installed, err
	}

	return installed,
		sh.RunV("npm", "install", "--prefix", path.Join(dst, n.versioned), n.Source+"@"+n.Version)
}
