package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const GitIgnoreExt = ".gitignore"
const GHRawUrl = "https://raw.githubusercontent.com/github/gitignore/master/"
const GHGitignoreRepo = "https://github.com/github/gitignore"

func listFiles() []string {
	var files []string

	c := colly.NewCollector(
		colly.AllowedDomains("github.com"),
	)

	c.OnHTML("table.files td.content a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasSuffix(href, GitIgnoreExt) {
			unscapedHref, _ := url.QueryUnescape(href)
			files = append(files, (path.Base(strings.TrimSuffix(unscapedHref, GitIgnoreExt))))
		}
	})

	c.Visit(GHGitignoreRepo)

	return files
}

func gitIgnoreExists(fileName string) (bool, error) {
	for _, f := range listFiles() {
		if f == fileName {
			return true, nil
		}
	}
	return false, fmt.Errorf(fileName + ".gitignore not found! Check the list with the -list option!")
}

func createUrl(fileName string) string {
	url := strings.Join([]string{GHRawUrl, fileName, GitIgnoreExt}, "")
	return url
}

func downloadFileContent(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return bodyBytes
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true, nil
	}
	return false, fmt.Errorf("File already exists")
}

func writeFileContent(content []byte) {
	writeFileError := ioutil.WriteFile(GitIgnoreExt, content, 0644)
	if writeFileError != nil {
		panic(writeFileError)
	}
}

func main() {
	list := flag.Bool("list", false, "list all .gitignore files")
	flag.Parse()

	if *list {
		fmt.Println(strings.Join(listFiles(), "\n"))
		os.Exit(0)
	}

	if len(os.Args) > 1 {
		arg := os.Args[1]

		_, fileExistsError := fileExists(GitIgnoreExt)
		if fileExistsError != nil {
			fmt.Println(fileExistsError)
			os.Exit(1)
		}

		_, gitIgnoreExistsError := gitIgnoreExists(arg)
		if gitIgnoreExistsError == nil {
			url := createUrl(arg)
			content := downloadFileContent(url)
			writeFileContent(content)
		} else {
			fmt.Println(gitIgnoreExistsError)
			os.Exit(1)
		}
	}
}
