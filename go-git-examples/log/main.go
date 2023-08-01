package main

import (
	"encoding/hex"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gogitexamples/util"
	"strings"
	"time"
)

func main() {
	myGitRepo := "./my_project"
	repo, err := git.PlainOpen(myGitRepo)
	util.CheckError(err)

	fmt.Println("all commits")
	cIter, err := repo.Log(&git.LogOptions{})
	util.CheckError(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	util.CheckError(err)

	fmt.Println()
	fmt.Println("all commits between since and until")
	since := time.Date(2021, 11, 30, 5, 47, 11, 0, time.UTC)
	until := time.Date(2021, 11, 30, 5, 47, 12, 0, time.UTC)
	cIter, err = repo.Log(&git.LogOptions{Since: &since, Until: &until})
	util.CheckError(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	util.CheckError(err)

	fmt.Println()
	fmt.Println("only commits in which a specific file was inserted/updated")
	file := "file1.md"
	cIter, err = repo.Log(&git.LogOptions{
		FileName: &file,
	})
	util.CheckError(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	util.CheckError(err)

	fmt.Println()
	fmt.Println("only commits that pass the path filter function")
	cIter, err = repo.Log(&git.LogOptions{
		PathFilter: func(s string) bool {
			return strings.HasSuffix(s, ".md")
		},
	})
	util.CheckError(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c.Hash, c.Message, c.Author.Email)
		return nil
	})
	util.CheckError(err)

	fmt.Println()
	fmt.Println("all commits starting from a specific git commit hash")
	hash, err := hex.DecodeString("0eb60866df254a9441ae4863808d631f7537a748")
	util.CheckError(err)

	var commitHash [20]byte
	copy(commitHash[:], hash)

	cIter, err = repo.Log(&git.LogOptions{
		From: commitHash,
	})
	util.CheckError(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c.Message)
		return nil
	})
	util.CheckError(err)
}
