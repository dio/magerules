package main

import (
	"fmt"

	"github.com/dio/magerules/example/pkg/version"
)

func main() {
	fmt.Println("foo", version.Commit, version.Version)
}
