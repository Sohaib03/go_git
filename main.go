package main

import (
	"os"

	"github.com/Sohaib03/go_git/cmd"
)

func main() {
	args := os.Args[1:]
	cmd.Call(args)
}
