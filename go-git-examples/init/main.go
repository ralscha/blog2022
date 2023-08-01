package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"gogitexamples/util"
)

func main() {
	repo, err := git.PlainInit("./my_project", false)
	util.CheckError(err)
	fmt.Println(repo)
}
