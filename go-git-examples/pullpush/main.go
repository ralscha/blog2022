package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"gogitexamples/util"
	"os"
	"time"
)

func main() {
	auth, err := ssh.NewPublicKeysFromFile("git", "/Users/myuser/.ssh/github", "")
	util.CheckError(err)

	r, err := git.PlainClone("./test", false, &git.CloneOptions{
		URL:      "github.com:ralscha/test.git",
		Progress: os.Stdout,
		Auth:     auth,
	})
	util.CheckError(err)

	w, err := r.Worktree()
	util.CheckError(err)

	err = w.Pull(&git.PullOptions{
		Auth:       auth,
		RemoteName: "origin",
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		util.CheckError(err)
	}

	readme := "./test/README.md"
	err = os.WriteFile(readme, []byte("Hello World 1"), 0644)
	util.CheckError(err)

	_, err = w.Add("README.md")
	util.CheckError(err)

	_, err = w.Commit("change README.md", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mr Author",
			Email: "author@email.com",
			When:  time.Now(),
		}})

	err = r.Push(&git.PushOptions{
		Auth: auth,
	})
	util.CheckError(err)
}
