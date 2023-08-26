package ko

import (
	"os"

	"github.com/dio/magerules/build"
	"github.com/dio/magerules/repo"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// ID only serves as an identifier internally
	ID string `yaml:",omitempty"`

	// Dir is the directory out of which the build should be triggered
	Dir string `yaml:",omitempty"`

	// Main points to the main package, or the source file with the main
	// function, in which case only the package will be used for the importpath
	Main string `yaml:",omitempty"`

	// Ldflags and Flags will be used for the Go build command line arguments
	Ldflags []string `yaml:",omitempty"`
	Flags   []string `yaml:",omitempty"`

	// Env allows setting environment variables for `go build`
	Env []string `yaml:",omitempty"`
}

func RemoveConfig() {
	_ = os.Remove(".ko.yaml")
}

type Entry struct {
	ID   string
	Main string
}

func WriteConfig(versionPkg string, entries ...Entry) error {
	mod, err := repo.GoMod()
	if err != nil {
		return err
	}
	var builds []Config
	for _, entry := range entries {
		buildFlags := build.StaticBuildFlags(mod, versionPkg)
		cfg := Config{
			ID:      entry.ID,
			Main:    entry.Main,
			Ldflags: buildFlags.LdFlags,
			Flags:   buildFlags.Flags,
		}
		builds = append(builds, cfg)
	}
	b, err := yaml.Marshal(map[string][]Config{"builds": builds})
	if err != nil {
		return err
	}
	return os.WriteFile(".ko.yaml", b, os.ModePerm)
}
