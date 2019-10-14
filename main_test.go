package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
  <title>GitHub Mock</title>
</head>
<body>
<table class="files">
  <td class="content">
    <a href="/toptal/gitignore/blob/master/templates/Go.gitignore">Go.gitignore</a>
  </td>
  <td class="content">
    <a href="/toptal/gitignore/blob/master/templates/Python.gitignore">Python.gitignore</a>
  </td>
</table>
</body>
</html>
		`))
	})

	return httptest.NewServer(mux)
}

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
	ts := newTestServer()
	if _, _, err := gitIgnoreExists(ts.URL, fileName); err != nil {
		t.Errorf("gitIgnoreExists(%s) was incorrect. Got: error to exist git ignore %s, want: nil error", fileName, fileName)
	}

	fileName = "gitIgnoreNotFound"
	if _, _, err := gitIgnoreExists(ts.URL, fileName); err == nil {
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
	ts := newTestServer()
	if files := listFiles(ts.URL); len(files) != 2 {
		t.Errorf("listFiles expected 2 files, but got %d files", len(files))
	}
}

func TestWriteFileContent(t *testing.T) {
	content, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		t.Errorf("Got: error %v, was: nil error to open .gitignore", err)
	}
	writeFileContent(content)
}
