package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

const (
	GitIgnoreExt    = ".gitignore"
	GHRawUrl        = "https://raw.githubusercontent.com/github/gitignore/master/"
	GHGitignoreRepo = "https://github.com/github/gitignore"
)

// listFiles ... describe what this function actually does (one line)
func listFiles() (files []string) {
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
	return
}

// gitIgnoreExists ... describe what this function actually does (one line)
func gitIgnoreExists(fileName string) (exist bool, err error) {
	for _, f := range listFiles() {
		if f == fileName {
			exist = true
			return
		}
	}
	err = fmt.Errorf(fileName + ".gitignore not found! Check the list with the -list option!")
	return
}

// createUrl ... describe what this function actually does (one line)
func createUrl(fileName string) (url string) {
	url = strings.Join([]string{GHRawUrl, fileName, GitIgnoreExt}, "")
	return
}

// downloadFileContent ... describe what this function actually does (one line)
func downloadFileContent(url string) (bodyBytes []byte) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// TODO: Never brush the error (_) and do not treat the error in the function, return to those who invoked
	bodyBytes, _ = ioutil.ReadAll(resp.Body)
	return
}

// fileExists ... describe what this function actually does (one line)
func fileExists(path string) (exit bool, err error) {
	if _, errStat := os.Stat(path); os.IsNotExist(errStat) {
		exit = true
		return
	}
	err = fmt.Errorf("File already exists")
	return
}

// writeFileContent ... describe what this function actually does (one line)
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
