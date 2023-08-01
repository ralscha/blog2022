package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gogitexamples/util"
	"os"
	"time"
)

func main() {
	myGitRepo := "./my_project"

	repo, err := git.PlainOpen(myGitRepo)
	util.CheckError(err)

	wt, err := repo.Worktree()
	util.CheckError(err)

	// create
	err = os.WriteFile(myGitRepo+"/file1.md", []byte("Hello World"), 0644)
	util.CheckError(err)

	_, err = wt.Add("file1.md")
	util.CheckError(err)

	hash, err := wt.Commit("create", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mr Author",
			Email: "author@email.com",
			When:  time.Now(),
		}})
	util.CheckError(err)
	fmt.Println("commit hash: ", hash)

	// update
	err = os.WriteFile(myGitRepo+"/file1.md", []byte("Hello Earth"), 0644)
	util.CheckError(err)

	err = wt.AddGlob(".")
	util.CheckError(err)

	_, err = wt.Commit("update", &git.CommitOptions{})
	util.CheckError(err)

	// delete
	err = os.Remove(myGitRepo + "/file1.md")
	util.CheckError(err)

	_, err = wt.Commit("delete", &git.CommitOptions{
		All: true,
	})

}
