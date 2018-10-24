package main

import "testing"

func TestCreateUrl(t *testing.T) {
	url := createUrl("Python")
	correctUrl := "https://raw.githubusercontent.com/dvcs/gitignore/master/templates/Python.gitignore"
	if url != correctUrl {
		t.Errorf("createUrl(%s) was incorrect. Got: %s, want: %s", "Python", url, correctUrl)
	}
}

func TestGitIgnoreExist(t *testing.T) {
	fileName := "Python"
	if _, _, err := gitIgnoreExists(fileName); err != nil {
		t.Errorf("[FAILED] - Exist gitignore for %s", fileName)
	}

	fileName = "gitIgnoreNotFound"
	if _, _, err := gitIgnoreExists(fileName); err == nil {
		t.Errorf("[FAILED] - Not exist gitignore for %s", fileName)
	}
}

func TestFileExists(t *testing.T) {
	gitIgnoreFile := ".gitignore"
	if _, err := fileExists(gitIgnoreFile); err == nil {
		t.Errorf("[FAILED] - File %s exist", gitIgnoreFile)
	}
}

func TestDownloadFileContent(t *testing.T) {
	url := "https://raw.githubusercontent.com/dvcs/gitignore/master/templates/Python.gitignore"
	if _, err := downloadFileContent(url); err != nil {
		t.Errorf("[FAILED] - Got error %v to download file of url %s", err, url)
	}
}

func TestListFiles(t *testing.T) {
	if files := listFiles(); len(files) == 0 {
		t.Error("[FAILED] - listFiles() returned an empty list")
	}
}
