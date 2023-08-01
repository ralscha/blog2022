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
	const gitRepo = "./my_project_branch"

	repo, err := git.PlainInit(gitRepo, false)
	util.CheckError(err)

	err = os.WriteFile(gitRepo+"/file1.md", []byte("Hello World"), 0644)
	util.CheckError(err)

	wt, err := repo.Worktree()
	util.CheckError(err)

	_, err = wt.Add(".")
	util.CheckError(err)

	_, err = wt.Commit("create files", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mr Author",
			Email: "author@email.com",
			When:  time.Now(),
		}})
	util.CheckError(err)

	// print current branch
	ref, err := repo.Head()
	util.CheckError(err)
	fmt.Println("current branch: ", ref.Name())

	// new branch
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("new_feature"), // or "refs/heads/new_feature"
		Create: true,
	})
	util.CheckError(err)

	err = os.WriteFile(gitRepo+"/file1.md", []byte("Hello Earth"), 0644)
	util.CheckError(err)
	err = os.WriteFile(gitRepo+"/branch_file.md", []byte("New Feature"), 0644)
	util.CheckError(err)

	_, err = wt.Add(".")
	util.CheckError(err)

	_, err = wt.Commit("implementing new feature", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mr Author",
			Email: "author@email.com",
			When:  time.Now(),
		}})
	util.CheckError(err)

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/master", // or omit. refs/heads/master is default
		Create: false,
	})
	util.CheckError(err)

	// list branches
	listBranches(repo)

	// delete branch
	err = repo.Storer.RemoveReference("refs/heads/new_feature")
	util.CheckError(err)

	fmt.Println("after delete branch")
	listBranches(repo)
}

func listBranches(repo *git.Repository) {
	branches, err := repo.Branches()
	util.CheckError(err)

	err = branches.ForEach(func(reference *plumbing.Reference) error {
		fmt.Println(reference.Name())
		return nil
	})
	util.CheckError(err)
}
