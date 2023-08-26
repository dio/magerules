package build

import (
	"fmt"

	"github.com/dio/magerules/repo"
	"github.com/dio/magerules/target"
	"github.com/magefile/mage/sh"
)

type BuildFlags struct {
	Flags   []string `yaml:"flags"`
	LdFlags []string `yaml:"ldflags"`
}

func StaticBuildFlags(mod, versionPkg string) BuildFlags {
	c, err := repo.ShortCommit()
	if err != nil {
		c = "unknown"
	}

	return BuildFlags{
		Flags: []string{"-tags", "-netgo"},
		LdFlags: []string{
			"-s",
			"-w",
			"-X", fmt.Sprintf("%s/%s.Commit=%s", mod, versionPkg, c),
			"-X", fmt.Sprintf("%s/%s.Version=%s", mod, versionPkg, repo.Version()),
			"-extldflags", "-static",
		},
	}
}

func StaticBuild(name, goos, goarch, versionPath string) error {
	x := target.Name(name)
	mod, err := repo.GoMod()
	if err != nil {
		return err
	}

	if versionPath == "" {
		versionPath = "pkg/version"
	}

	buildFlags := StaticBuildFlags(mod, versionPath)
	buildArgs := []string{"build"}
	buildArgs = append(buildArgs, buildFlags.Flags...)
	buildArgs = append(buildArgs, "-ldflags")
	buildArgs = append(buildArgs, buildFlags.LdFlags...)

	output := fmt.Sprintf("build/%s_%s_%s_%s/%s", x.Key(), repo.Version(), goos, goarch, x.Name())
	buildArgs = append(buildArgs, output, mod+"/"+x.Key())

	env := map[string]string{
		"GOOS":        goos,
		"GOARCH":      goarch,
		"CGO_ENABLED": "0",
	}
	return sh.RunWith(env, "go", buildArgs...)
}
