package tools

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dio/magerules/installable"
	"github.com/dio/magerules/repo"
	"github.com/magefile/mage/sh"
)

func LoadDefault() (*Tools, error) {
	return Load("magetools")
}

func Load(dir string) (*Tools, error) {
	data, err := os.ReadFile(".magetools.yaml")
	if err != nil {
		return nil, err
	}
	installables, err := installable.LoadInstallables(data)
	if err != nil {
		return nil, err
	}

	if !path.IsAbs(dir) {
		topLevel, _ := repo.ShowToplevel()
		dir = path.Join(topLevel, dir)
	}

	return &Tools{Installables: installables, Dir: dir}, nil
}

type Tools struct {
	Dir          string
	Installables installable.Installables
}

func (t *Tools) Run(name string, args ...string) error {
	return t.RunWith(RunWithOptions{}, name, args...)
}

// RunWithOptions provides options when running a tool.
type RunWithOptions struct {
	Deps []string
	Env  map[string]string
}

func (t *Tools) RunWith(opt RunWithOptions, name string, args ...string) error {
	opt.Deps = append(opt.Deps, name)
	p, err := t.Install(opt.Deps...)
	if err != nil {
		return err
	}
	os.Setenv("PATH", p+":"+os.Getenv("PATH"))

	info, err := t.getInstallableInfo(name)
	if err != nil {
		return err
	}
	return sh.RunWithV(opt.Env, info.Binary, args...)
}

func (t *Tools) Install(names ...string) (string, error) {
	var paths []string
	for _, name := range names {
		info, err := t.getInstallableInfo(name)
		if err != nil {
			return strings.Join(paths, ":"), nil
		}

		for _, i := range info.Installers {
			p, err := i.Install(t.Dir)
			paths = append(paths, p)
			if err != nil {
				_ = os.RemoveAll(installable.GetInstalledBase(p))
				return strings.Join(paths, ":"), err
			}
		}
	}
	return strings.Join(paths, ":"), nil
}

func (t *Tools) InstallAll() error {
	var names []string
	for name := range t.Installables {
		names = append(names, name)
	}
	if _, err := t.Install(names...); err != nil {
		return err
	}
	return nil
}

func (t *Tools) getInstallableInfo(name string) (installable.InstallabeInfo, error) {
	info := installable.InstallabeInfo{Key: name, Binary: name}
	if strings.Contains(name, ":") {
		parts := strings.Split(name, ":")
		if len(parts) != 2 {
			return info, fmt.Errorf("name: %s %w", name, installable.ErrInvalid)
		}
		info.Key = parts[0]
		info.Binary = parts[1]
	}

	i, ok := t.Installables[info.Key]
	if !ok {
		return info, fmt.Errorf("unknown name: %s %w", name, installable.ErrInvalid)
	}
	if i.Runtime() != nil {
		info.Installers = append(info.Installers, i.Runtime())
	}
	info.Installers = append(info.Installers, i)
	return info, os.MkdirAll(t.Dir, os.ModePerm)
}
