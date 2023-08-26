package build

import (
	"fmt"

	"github.com/dio/magerules/repo"
	"github.com/dio/magerules/target"
	"github.com/magefile/mage/sh"
)

func StaticBuild(name, goos, goarch string) error {
	x := target.Name(name)
	mod, err := repo.GoMod()
	if err != nil {
		return err
	}

	ldflags, err := StaticLdflags(mod)
	if err != nil {
		return err
	}

	env := map[string]string{
		"GOOS":   goos,
		"GOARCH": goarch,
	}
	output := fmt.Sprintf("build/%s_%s_%s_%s/%s", x.Key(), repo.Version(), goos, goarch, x.Name())
	return sh.RunWith(env, "go", "build", "-ldflags", ldflags, "-o", output, mod+"/"+x.Key())
}

func StaticLdflags(mod string) (string, error) {
	c, err := repo.ShortCommit()
	if err != nil {
		return "", err
	}
	// It is advisable to have a pkg/version on your project.
	return fmt.Sprintf("-s -w -X %s/pkg/version.Commit=%s -X %s/pkg/version.Version=%s", mod, c, mod, repo.Version()), nil
}
