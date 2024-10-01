package pkg

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strconv"
)

type GitObject interface {
	fmt() string
	Serialize() []byte
	Deserialize(data []byte) error
}

type GitBlob struct {
	data []byte
}

func (b *GitBlob) fmt() string {
	return "blob"
}

func (b *GitBlob) Serialize() []byte {
	return b.data
}

func (b *GitBlob) Deserialize(data []byte) error {
	b.data = data
	return nil
}

func ObjectFind(repo *Repository, name string, p string, follow bool) string {
	return name
}

func ObjectRead(repo *Repository, sha string) (GitObject, error) {
	filePath := repoPath(repo, "objects", sha[:2], sha[2:])

	fmt.Println(filePath)

	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("Object not found: %s", sha)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening object file: %s", filePath)
	}
	defer file.Close()

	compressedData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading object file: %s", filePath)
	}

	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("Error creating zlib reader: %s", filePath)
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Error reading zlib reader: %s", filePath)
	}

	x := bytes.IndexByte(raw, ' ')
	if x == -1 {
		return nil, fmt.Errorf("Invalid object file: %s", filePath)
	}
	fmtType := string(raw[:x])
	y := bytes.IndexByte(raw, 0)
	if y == -1 {
		return nil, fmt.Errorf("Invalid object file: %s", filePath)
	}
	// y += x
	size, err := strconv.Atoi(string(raw[x+1 : y]))
	if err != nil {
		return nil, fmt.Errorf("Invalid object file: %s", filePath)
	}
	fmt.Println("File info : ", fmtType, size)

	switch fmtType {
	case "commit":
		break
	case "tree":
		break
	case "blob":
		return &GitBlob{data: raw[y+1:]}, nil
		break
	case "tag":
		break
	default:
		return nil, fmt.Errorf("Unknown object type: %s", fmtType)
	}
	return nil, nil
}

func ObjectWrite(repo *Repository, obj GitObject) (string, error) {
	data := obj.Serialize()
	header := obj.fmt() + " " + strconv.Itoa(len(data))
	byte_data := []byte(header)
	byte_data = append(byte_data, 0)
	byte_data = append(byte_data, data...)

	sha := fmt.Sprintf("%x", sha1.Sum(byte_data))

	if repo != nil {
		rpath, err := repoFile(repo, true, "objects", sha[:2], sha[2:])
		fmt.Println("rpath: ", rpath)
		if err == nil {
			file, err := os.Create(rpath)
			// file, err := os.Open(rpath)
			defer file.Close()
			if err != nil {
				fmt.Print("Error creating object file: %s", rpath)
				return "", fmt.Errorf("Error opening object file: %s", rpath)
			} else {
				fmt.Println("File Created")
				writer := zlib.NewWriter(file)
				defer writer.Close()
				n, err := writer.Write(byte_data)
				if err != nil {
					return "", fmt.Errorf("Error writing object file: %s", rpath)
				} else if n != len(byte_data) {
					return "", fmt.Errorf("Written %d / %d bytes: %s", n, len(byte_data), rpath)
				}
				fmt.Println("Written to ", file.Name())
			}
		}
	}
	return sha, nil
}
