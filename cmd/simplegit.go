package cmd

import (
	"fmt"

	"github.com/Sohaib03/go_git/pkg"
)

func Call(args []string) {
	fmt.Println("Argument count:", len(args))
	fmt.Printf("Arguments: %v\n", args)

	if len(args) == 0 {
		// show usage
		fmt.Println("Usage: simplegit <command> [<args>]")
		return
	}

	_, err := pkg.GetRepository(".", false)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch args[0] {
	case "init":
		fmt.Println("Initializing repository...")
	case "commit":
		fmt.Println("Committing changes...")
	case "status":
		fmt.Println("Checking status...")
	default:
		fmt.Println("Unknown command:", args[0])
	}
}
