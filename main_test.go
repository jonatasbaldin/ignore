package main

import (
	"io/ioutil"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestCreateUrl(t *testing.T) {
	url := createUrl("Python")
	correctUrl := "https://raw.githubusercontent.com/dvcs/gitignore/master/templates/Python.gitignore"
	if url != correctUrl {
		t.Errorf("createUrl(%s) was incorrect. Got: %s, want: %s", "Python", url, correctUrl)
	}
}

func TestDownloadFileContent(t *testing.T) {
	url := "http://foo.com"
	defer gock.Off()

	gock.New(url).
		Get("/").
		Reply(200).
		BodyString("Foo")

	resp, err := downloadFileContent(url)
	if string(resp) != "Foo" {
		t.Errorf(
			"downloadFileContent(%s) was incorrect. Got: %s, want: %s",
			url,
			string(resp),
			"Foo",
		)
	}
	if err != nil {
		t.Errorf("[FAILED] - Got error %v to download file of url %s", err, url)
	}
}

func TestGitIgnoreExist(t *testing.T) {
	fileName := "Python"
	if _, _, err := gitIgnoreExists(fileName); err != nil {
		t.Errorf("[FAILED] - Exist gitignore for %s", fileName)
		t.Errorf("gitIgnoreExists(%s) was incorrect. Got: error to exist git ignore %s, want: nil error", fileName, fileName)
	}

	fileName = "gitIgnoreNotFound"
	if _, _, err := gitIgnoreExists(fileName); err == nil {
		t.Errorf("gitIgnoreExists(%s) was incorrect. Got: nil error, want: error to exist file %s", fileName, fileName)
	}
}

func TestFileExists(t *testing.T) {
	gitIgnoreFile := ".gitignore"
	if _, err := fileExists(gitIgnoreFile); err == nil {
		t.Errorf("fileExists(%s) was incorrect. Got: nil error, want: error file %s already exist", gitIgnoreFile, gitIgnoreFile)
	}
}

func TestListFiles(t *testing.T) {
	if files := listFiles(); len(files) == 0 {
		t.Error("listFiles() was incorrect. Got empty list")
	}
}

func TestWriteFileContent(t *testing.T) {
	content, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		t.Errorf("Got: error %v, was: nil error to open .gitignore", err)
	}
	writeFileContent(content)
}
