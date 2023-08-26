//go:build mage
// +build mage

package main

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/dio/magerules/build"
	"github.com/dio/magerules/dist"
	"github.com/dio/magerules/repo"
	"github.com/dio/magerules/tools"
)

func Build(targets string) error {
	names := strings.Split(targets, ",")
	for _, name := range names {
		if err := build.StaticBuild(name, runtime.GOOS, runtime.GOARCH, "pkg/version"); err != nil {
			return err
		}
	}
	return nil
}

func Dist(targets string) error {
	names := strings.Split(targets, ",")
	for _, name := range names {
		name := name
		if err := dist.Create(name, "pkg/version", build.StaticBuild); err != nil {
			return err
		}
	}
	return nil
}

func Clean() {
	_ = os.RemoveAll("build")
	_ = os.RemoveAll("dist")
}

func Generate() error {
	toolbox, err := tools.LoadDefault()
	if err != nil {
		return err
	}

	return toolbox.RunWith(tools.RunWithOptions{Deps: []string{
		"protoc-gen-go",
		"protoc-gen-connect-go",
		"protoc-gen-es",         // this requires "node", see: .magetools.yaml.
		"protoc-gen-connect-es", // this requires "node", see: .magetools.yaml.
	}}, "buf", "generate")
}

func Bundle(name string) error {
	// TODO(dio): Add a dep or a flag to force to re-generating api.
	toolbox, err := tools.Load(topLevel())
	if err != nil {
		return err
	}

	if err = toolbox.Run("node:corepack", "enable"); err != nil {
		return err
	}

	if err = toolbox.Run("node:yarn", "install", "--immutable"); err != nil {
		return err
	}

	if !strings.HasPrefix(name, "@ui") {
		name = path.Join("@ui", name)
	}
	return toolbox.Run("node:yarn", "nx", "build", name)
}

func Tools() error {
	toolbox, err := tools.Load(topLevel())
	if err != nil {
		return err
	}
	return toolbox.InstallAll()
}

// Using defaut, we seek for --show-toplevel. However, this is in example.
func topLevel() string {
	p, _ := repo.ShowToplevel()
	return path.Join(p, "example", "magetools")
}
