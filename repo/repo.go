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
	list, err := sh.Output("go", "list", "-m")
	if err != nil {
		return "", err
	}
	// The "list" contains multiple modules separated by "\n". We are interested in the first entry.
	mods := strings.SplitN(list, "\n", 2)
	return mods[0], nil
}
