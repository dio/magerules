package image

import (
	"github.com/dio/magerules/image/ko"
	"github.com/dio/magerules/tools"
)

type PushOptions struct {
	Entries    []ko.Entry
	Tag        string
	VersionPkg string
	Toolbox    *tools.Tools
	Platforms  []string
	Push       bool
	Tarball    string
	Registry   string
}

func Push(opt PushOptions) error {
	// if err := ko.WriteConfig(opt.VersionPkg, opt.Entries...); err != nil {
	// 	return err
	// }
	// defer ko.RemoveConfig()

	// mod, err := repo.GoMod()
	// if err != nil {
	// 	return err
	// }

	// for _, entry := range opt.Entries {
	// 	// Build image names, since we do bare

	// }

	// return opt.Toolbox.RunWith(tools.RunWithOptions{
	// 	Env: map[string]string{
	// 		// "KO_DOCKER_REPO": opt.imageFqn,
	// 	},
	// },
	// 	"ko", "build", path.Join(mod, name),
	// 	"--tags", tag,
	// 	"--platform", "linux/amd64",
	// 	"--platform", "linux/arm64",
	// 	"--bare",
	// 	"--push=false",
	// 	"--tarball",
	// 	"/tmp/a.tar.gz",
	// )
	return nil
}
