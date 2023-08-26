package dist

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/dio/magerules/repo"
	"github.com/dio/magerules/target"
	"github.com/magefile/mage/sh"
)

func Create(name string, build func(string, string, string) error) error {
	var archives []string
	x := target.Name(name)
	for _, goos := range []string{"linux", "darwin"} {
		goos := goos
		for _, goarch := range []string{"arm64", "amd64"} {
			goarch := goarch
			err := build(name, goos, goarch)
			if err != nil {
				return err
			}
			targz, err := createTargz(x, goos, goarch)
			if err != nil {
				return err
			}
			archives = append(archives, targz)
		}
	}
	return createChecksums(x.Key(), archives)
}

func createTargz(name target.Name, goos, goarch string) (string, error) {
	dir := fmt.Sprintf("%s_%s_%s_%s", name.Key(), repo.Version(), goos, goarch)
	src := path.Join("build", dir)
	dst := path.Join("dist", dir)

	err := os.MkdirAll("dist", os.ModePerm)
	if err != nil {
		return "", err
	}

	targz := path.Join(dst + ".tar.gz")
	err = sh.RunV("tar", "-C", src, "-cpzf", targz, name.Name())
	if err != nil {
		return "", err
	}

	return targz, nil
}

func createChecksums(name string, archives []string) error {
	args := []string{}
	if runtime.GOOS == "darwin" {
		args = append(args, "-a", "256")
	}
	args = append(args, archives...)
	out, err := sh.Output("shasum", args...)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join("dist", fmt.Sprintf("%s_%s_checksums.txt", name, repo.Version())), []byte(out), os.ModePerm)
}
