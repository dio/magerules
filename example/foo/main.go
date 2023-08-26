package main

import (
	"fmt"

	"github.com/dio/magerules/example/pkg/version"
)

var commit string

func main() {
	fmt.Println("haha", version.Commit, version.Version, commit)
}
