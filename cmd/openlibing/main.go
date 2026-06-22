package main

import (
	"fmt"
	"io/fs"

	embedded "github.com/openlibing/openlibing-cli/embedded"
)

func main() {
	entries, _ := fs.ReadDir(embedded.SPCs, "spc")
	for _, e := range entries {
		fmt.Println(e.Name())
	}
}
