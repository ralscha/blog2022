package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"gogitexamples/util"
	"os"
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

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "example",
		URLs: []string{"https://github.com/ralscha/test2.git"},
	})
	util.CheckError(err)

	list, err := r.Remotes()
	util.CheckError(err)

	for _, r := range list {
		fmt.Println(r)
	}

	err = r.DeleteRemote("example")
	util.CheckError(err)

	fmt.Println("After delete")
	list, err = r.Remotes()
	util.CheckError(err)

	for _, r := range list {
		fmt.Println(r)
	}

}
