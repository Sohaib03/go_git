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

	// repo, err := pkg.GetRepository(".", false)

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// _, err = pkg.ObjectRead(repo, "f73c96d8a7de0a95b3be41d186feef265c63a980")
	// fmt.Println(err)

	switch args[0] {
	case "init":
		fmt.Println("Initializing repository...")
		if len(args) > 1 {
			pkg.RepoCreate(args[1])
		} else {
			pkg.RepoCreate(".")
		}
	case "commit":
		fmt.Println("Committing changes...")
	case "status":
		fmt.Println("Checking status...")
		pkg.RepoFind(".", true)

	case "hash":
		fmt.Println("Hashing file...")
		pkg.CmdHashFile(args)

	case "cat":
		fmt.Println("Reading file...")
		pkg.CmdCatFile(args)
	default:
		fmt.Println("Unknown command:", args[0])
	}
}
