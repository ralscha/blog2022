package main

import (
	"github.com/go-git/go-git/v5"
	"gogitexamples/util"
	"os"
)

func main() {
	_, err := git.PlainClone("./go-git", true, &git.CloneOptions{
		URL:      "https://github.com/go-git/go-git",
		Progress: os.Stdout,
	})

	util.CheckError(err)
}
