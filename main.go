package main

import (
	"os"

	"github.com/scttfrdmn/gh-milestone/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
