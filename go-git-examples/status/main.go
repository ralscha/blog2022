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
	const gitRepo = "./my_project_status"

	repo, err := git.PlainInit(gitRepo, false)
	util.CheckError(err)

	wt, err := repo.Worktree()
	util.CheckError(err)

	err = os.WriteFile(gitRepo+"/ignored_file.md", []byte("Ignore me"), 0644)
	util.CheckError(err)

	err = os.WriteFile(gitRepo+"/.gitignore", []byte("ignored_file.md"), 0644)
	util.CheckError(err)

	err = os.WriteFile(gitRepo+"/file1.md", []byte("Hello World 1"), 0644)
	util.CheckError(err)
	err = os.WriteFile(gitRepo+"/file2.md", []byte("Hello World 2"), 0644)
	util.CheckError(err)
	err = os.WriteFile(gitRepo+"/file3.md", []byte("Hello World 3"), 0644)
	util.CheckError(err)
	err = os.WriteFile(gitRepo+"/fileOldName.md", []byte("File with old name"), 0644)
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

	err = os.WriteFile(gitRepo+"/file4.md", []byte("Hello World 4"), 0644)
	util.CheckError(err)
	err = os.WriteFile(gitRepo+"/file5.md", []byte("Hello World 5"), 0644)
	util.CheckError(err)

	_, err = wt.Add("file5.md")

	err = os.WriteFile(gitRepo+"/file2.md", []byte("Hello Earth 2"), 0644)
	util.CheckError(err)

	_, err = wt.Add("file2.md")

	err = os.Remove(gitRepo + "/file1.md")
	util.CheckError(err)
	_, err = wt.Add("file1.md")

	err = os.Remove(gitRepo + "/file3.md")
	util.CheckError(err)

	_, err = wt.Move("fileOldName.md", "fileNewName.md")
	util.CheckError(err)

	status, err := wt.Status()
	util.CheckError(err)

	fmt.Printf("is clean %t\n", status.IsClean())

	for file, status := range status {
		fmt.Print(file + ": ")
		fmt.Print("Staging: " + toHumanReadable(status.Staging))
		fmt.Print("  ")
		fmt.Print("Worktree: " + toHumanReadable(status.Worktree))
		fmt.Println()
	}

}

func toHumanReadable(statusCode git.StatusCode) string {
	switch statusCode {
	case git.Unmodified:
		return "unmodified"
	case git.Untracked:
		return "untracked"
	case git.Modified:
		return "modified"
	case git.Added:
		return "added"
	case git.Deleted:
		return "deleted"
	case git.Renamed:
		return "renamed"
	case git.Copied:
		return "copied"
	case git.UpdatedButUnmerged:
		return "updated but unmerged"
	}
	return ""
}
