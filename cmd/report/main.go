package main

import (
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
)

const (
	RepoPath = "/Users/tuanvuong/Workspace/nordotstrictrss/nordot-feedreader-fvc"
)

func main() {
	repo, err := git.PlainOpen(RepoPath)
	if err != nil {
		panic(err)
	}

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		commit, err := iter.Next()
		if err != nil {
			break
		}
		msg := fmt.Sprintf("%s", commit.String())
		fmt.Println(msg)
	}
}
