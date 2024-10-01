package pkg

import (
	"fmt"
	"io"
	"os"
)

func CmdCatFile(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: simplegit cat <object>")
		return
	}

	repo := RepoFind(".", true)

	err := CatFile(repo, args[1], "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func CatFile(repo *Repository, object string, p string) error {
	obj, err := ObjectRead(repo, ObjectFind(repo, object, p, true))
	if err != nil {
		return err
	}
	fmt.Println("Object :", string(obj.Serialize()))
	return nil
}

func CmdHashFile(args []string) {
	// hash-object [-t TYPE] FILE
	if len(args) < 3 {
		fmt.Println("Usage: simplegit hash-object [-t TYPE] FILE")
		return
	}

	repo := RepoFind(".", true)

	fd := args[2]
	objectType := args[1]

	file, err := os.Open(fd)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	s, err := HashFile(repo, file, objectType)
	if err != nil {
		fmt.Println("Error hashing file:", err)
		return
	}
	fmt.Println(s)

}

func HashFile(repo *Repository, file *os.File, objectType string) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error reading file: %s", file.Name())
	}
	switch objectType {
	case "blob":
		blob := &GitBlob{data}
		return ObjectWrite(repo, blob)
		break
	default:
		return "", fmt.Errorf("Unknown object type: %s", objectType)
	}

	return "", nil
}
