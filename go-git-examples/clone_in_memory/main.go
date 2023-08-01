package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"gogitexamples/util"
)

func main() {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/go-git/go-git",
	})
	util.CheckError(err)

	ref, err := r.Head()
	util.CheckError(err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	util.CheckError(err)

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	util.CheckError(err)
}
