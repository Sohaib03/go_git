package pkg

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	// "github.com/bigkevmcd/go-configparser"
	"gopkg.in/ini.v1"
)

type Repository struct {
	worktree string
	gitdir   string
	conf     *ini.File
}

func GetRepository(path string, force bool) (*Repository, error) {
	repo := &Repository{}
	return repo, repo.init(path, force)
}

func repoPath(r *Repository, paths ...string) string {
	paths = append([]string{r.gitdir}, paths...)
	return path.Join(paths...)
}

func repoFile(r *Repository, mkdir bool, paths ...string) (string, error) {
	if _, err := repoDir(r, mkdir, paths[:len(paths)-1]...); err != nil {
		return "", err
	}
	// Return full path
	return repoPath(r, paths...), nil
}

func repoDir(r *Repository, mkdir bool, paths ...string) (string, error) {
	curPath := repoPath(r, paths...)

	// check if path exists
	if _, err := os.Stat(curPath); err == nil {
		// check if path is a directory
		if stat, err := os.Stat(curPath); err == nil && stat.IsDir() {
			return curPath, nil
		} else {
			return "", fmt.Errorf("Not a directory %s", curPath)
		}
	}

	if mkdir {
		if err := os.MkdirAll(curPath, 0755); err != nil {
			return "", fmt.Errorf("Error creating directory %s", curPath)
		}
		return curPath, nil
	} else {
		return "", fmt.Errorf("Directory does not exist %s", curPath)
	}

}

// Initializes a Git repository at rpath; force skips validation checks.
func (r *Repository) init(rpath string, force bool) error {
	r.worktree = rpath
	r.gitdir = path.Join(rpath, ".git")
	// check if the directory exists
	if !force {
		if stat, err := os.Stat(r.gitdir); err != nil || !stat.IsDir() {
			return fmt.Errorf("Invalid Git Repository %s\n", r.gitdir)
		}
	}
	conf, err := ini.Load(".git/config")
	// conf, err := configparser.NewConfigParserFromFile(".git/config")
	if err != nil {
		return fmt.Errorf("Error reading config file")
	}
	r.conf = conf

	if !force {
		data, err := conf.Section("core").GetKey("repositoryformatversion")
		if err != nil {
			return fmt.Errorf("Error reading repositoryformatversion")
		}
		version := data.String()
		// TODO : Check if the version is supported
		if version != "0" {
			return fmt.Errorf("Unsupported repositoryformatversion %s", version)
		}
	}

	fmt.Println("Repository Found")
	return nil
}

func RepoCreate(rpath string) {
	fmt.Println("Creating repository at", rpath)
	repo := &Repository{}
	repo.init(rpath, true)

	// check if the directory exists
	if stat, err := os.Stat(repo.worktree); err == nil {
		// check is dir
		// fmt.Println("HERE Error:", err)
		if !stat.IsDir() {
			fmt.Println("Error: Not a directory", repo.worktree)
			return
		}

		// check if gitdir is empty
		if files, err := os.ReadDir(repo.gitdir); err == nil && len(files) > 0 {
			fmt.Println("Error: Not an empty directory", repo.worktree)
			return
		}

	} else {
		// create the directory
		if err := os.MkdirAll(repo.worktree, 0755); err != nil {
			fmt.Println("Error: Could not create directory", repo.worktree)
			return
		}
	}

	if _, err := repoDir(repo, true, "branches"); err != nil {
		fmt.Println("Error in branches:", err)
		return
	}
	if _, err := repoDir(repo, true, "objects"); err != nil {
		fmt.Println("Error in objects:", err)
		return
	}
	if _, err := repoDir(repo, true, "refs", "tags"); err != nil {
		fmt.Println("Error in refs/tags:", err)
		return
	}
	if _, err := repoDir(repo, true, "refs", "heads"); err != nil {
		fmt.Println("Error in refs/heads:", err)
		return
	}

	// write to repoFile(repo, false, "description")
	filepath, err := repoFile(repo, false, "description")
	if err != nil {
		fmt.Println("Error in description:", err)
		return
	}
	if err := os.WriteFile(filepath,
		[]byte("Unnamed repository; edit this file 'description' to name the repository.\n"),
		0644); err != nil {
		fmt.Println("Error writing to description:", err)
		return
	}

	// write to repoFile(repo, false, "HEAD")
	filepath, err = repoFile(repo, false, "HEAD")
	if err != nil {
		fmt.Println("Error in HEAD:", err)
		return
	}
	if err := os.WriteFile(filepath, []byte("ref: refs/heads/master\n"), 0644); err != nil {
		fmt.Println("Error writing to HEAD:", err)
		return
	}

	// write to repoFile(repo, false, "config")
	filepath, err = repoFile(repo, false, "config")
	if err != nil {
		fmt.Println("Error in config:", err)
		return
	}

	err = writeDefaultConfig(filepath)
	if err != nil {
		fmt.Println("Error writing to config:", err)
		return
	}

}

func writeDefaultConfig(filepath string) error {
	config := ini.Empty()
	config.Section("core").Key("repositoryformatversion").SetValue("0")
	config.Section("core").Key("filemode").SetValue("false")
	config.Section("core").Key("bare").SetValue("false")
	return config.SaveTo(filepath)
}

func RepoFind(cpath string, required bool) *Repository {
	realPath, err := filepath.Abs(cpath)
	if err != nil {
		fmt.Println("Error getting absolute path")
		return nil
	}
	for {
		gitPath := path.Join(realPath, ".git")
		if stat, err := os.Stat(gitPath); err == nil && stat.IsDir() {
			fmt.Println("Repository Found at ", realPath)
			repo := &Repository{}
			repo.init(realPath, false)
			return repo
		}
		fmt.Println("Not Found at ", realPath)
		if realPath == "/" || realPath == "." || realPath == filepath.Dir(realPath) {
			if required {
				fmt.Println("No repository found")
				return nil
			}
			return nil
		}
		realPath = filepath.Dir(realPath)
	}
}
