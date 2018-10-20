package main

import "testing"

func TestCreateUrl(t *testing.T) {
	url := createUrl("Python")
	correctUrl := "https://raw.githubusercontent.com/dvcs/gitignore/master/templates/Python.gitignore"
	if url != correctUrl {
		t.Errorf("createUrl(%s) was incorrect. Got: %s, want: %s", "Python", url, correctUrl)
	}
}
