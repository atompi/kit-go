package git

import (
	"github.com/go-git/go-git/v5"
)

func Clone(url string, workDir string) (r *git.Repository, err error) {
	o := &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.NoRecurseSubmodules,
	}
	r, err = git.PlainClone(workDir, false, o)
	if err != nil {
		return
	}
	return
}

func RemoveRemote(r *git.Repository, remoteName string) (err error) {
	err = r.DeleteRemote(remoteName)
	return
}

func Open(path string) (r *git.Repository, err error) {
	r, err = git.PlainOpen(path)
	return
}
