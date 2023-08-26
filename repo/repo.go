package repo

import (
	"os"
	"strings"

	"github.com/magefile/mage/sh"
)

func ShowToplevel() (string, error) {
	return sh.Output("git", "rev-parse", "--show-toplevel")
}

func ShortCommit() (string, error) {
	return sh.Output("git", "rev-parse", "--short", "HEAD")
}

func Version() string {
	ver := os.Getenv("VERSION")
	if ver == "" {
		return "dev"
	}
	return ver
}

func GoMod() (string, error) {
	list, err := sh.Output("go", "mod", "edit", "-print")
	if err != nil {
		return "", err
	}
	// The "list" contains multiple modules separated by "\n". We are interested in the first entry.
	mods := strings.SplitN(list, "\n", 2)
	mod := strings.SplitN(mods[0], " ", 2)
	return mod[1], nil
}
