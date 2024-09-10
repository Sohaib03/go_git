package pkg

import (
	"fmt"
	"os"
	"path"

	"github.com/bigkevmcd/go-configparser"
)

type Repository struct {
	worktree string
	gitdir   string
	conf     *configparser.ConfigParser
}

func GetRepository(path string, force bool) (*Repository, error) {
	repo := &Repository{}
	return repo, repo.init(path, force)
}

func (r *Repository) init(rpath string, force bool) error {
	r.worktree = rpath
	r.gitdir = path.Join(rpath, ".git")
	// check if the directory exists
	if !force {
		if stat, err := os.Stat(r.gitdir); err != nil || !stat.IsDir() {
			return fmt.Errorf("Invalid Git Repository", r.gitdir)
		}
	}
	conf, err := configparser.NewConfigParserFromFile(".git/config")
	if err != nil {
		return fmt.Errorf("Error reading config file")
	}
	r.conf = conf

	if !force {
		_, err := conf.Get("core", "repositoryformatversion")
		if err != nil {
			return fmt.Errorf("Error reading repositoryformatversion")
		}
		// TODO : Check if the version is supported
		// if vers != "0" {
		// 	return fmt.Errorf("Unsupported repositoryformatversion %s", vers)
		// }
	}

	fmt.Println("Repository Found")
	return nil
}
