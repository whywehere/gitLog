package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"testing"
)

func TestGoGit(t *testing.T) {
	r, err := git.PlainOpen("C:\\Users\\19406\\Desktop\\go\\tta\\cinx")
	if err != nil {
		t.Fatal(err)
	}
	head, err := r.Head()
	if err != nil {
		t.Fatal(err)
	}
	cIter, _ := r.Log(&git.LogOptions{From: head.Hash()})
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Printf("email: %v\n", c.Author.Email)
		return nil
	})

}
