package main

import (
	"archive/zip"
	"compress/flate"
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	token := os.Getenv("GITHUB_BACKUP_TOKEN")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{
		Visibility:  "all",
		Affiliation: "owner",
		Sort:        "full_name",
		Direction:   "asc",
		ListOptions: github.ListOptions{PerPage: 25},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(context.Background(), "", opt)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	auth, err := ssh.NewPublicKeysFromFile("git", "/Users/myuser/.ssh/github", "")
	if err != nil {
		log.Fatal(err)
	}

	backupDir := "./backup_github"

	for _, repo := range allRepos {
		_, err := git.PlainClone(backupDir+"/"+*repo.FullName, true, &git.CloneOptions{
			Auth:     auth,
			URL:      *repo.SSHURL,
			Progress: os.Stdout,
		})
		if err != nil && errors.Is(err, git.ErrRepositoryAlreadyExists) {
			fmt.Println("pulling", *repo.FullName)
			r, err := git.PlainOpen("./backup_github/" + *repo.FullName)
			if err != nil {
				log.Fatal(err)
			}
			err = r.Fetch(&git.FetchOptions{
				Auth: auth,
			})
			if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
				log.Fatal(err)
			}
		}
	}

	zipFile := backupDir + ".zip"
	err = zipit(backupDir, zipFile)
	if err != nil {
		log.Fatal(err)
	}
}

func zipit(backupDir, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	defer zipWriter.Close()

	return filepath.Walk(backupDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		fileName := filepath.ToSlash(filePath)
		index := strings.Index(fileName, "/")
		if index != -1 {
			fileName = fileName[index+1:]
		}

		writer, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

}
