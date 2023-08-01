package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gogitexamples/util"
	"os"
	"time"
)

func main() {
	myGitRepo := "./my_project_tag"

	repo, err := git.PlainInit(myGitRepo, false)
	util.CheckError(err)

	wt, err := repo.Worktree()
	util.CheckError(err)

	err = os.WriteFile(myGitRepo+"/file1.md", []byte("Hello World 1"), 0644)
	util.CheckError(err)

	_, err = wt.Add(".")
	util.CheckError(err)

	hash, err := wt.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mr Author",
			Email: "author@email.com",
			When:  time.Now(),
		}})
	util.CheckError(err)

	_, err = repo.CreateTag("v1.0", hash, &git.CreateTagOptions{
		Message: "version 1.0",
	})
	util.CheckError(err)

	err = os.WriteFile(myGitRepo+"/file2.md", []byte("Hello World 2"), 0644)
	util.CheckError(err)

	_, err = wt.Add(".")
	util.CheckError(err)

	hash, err = wt.Commit("new feature", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mr Author",
			Email: "author@email.com",
			When:  time.Now(),
		}})
	util.CheckError(err)

	_, err = repo.CreateTag("v1.0.1", hash, &git.CreateTagOptions{
		Message: "version 1.0.1",
	})
	util.CheckError(err)

	err = repo.DeleteTag("v1.0.1")
	util.CheckError(err)

	_, err = repo.CreateTag("v1.1", hash, &git.CreateTagOptions{
		Message: "version 1.1",
	})
	util.CheckError(err)

	tags, err := repo.Tags()
	util.CheckError(err)

	err = tags.ForEach(func(reference *plumbing.Reference) error {
		fmt.Println(reference.Name())
		return nil
	})
	util.CheckError(err)
}
