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

// return a slice of strings containing all the .gitignore file names from GitHub
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

// checks if the specified file name exists in the .gitignore list of files from GitHub
func gitIgnoreExists(fileName string) (exist bool, name string, err error) {
	for _, f := range listFiles() {
		if strings.EqualFold(f, fileName) {
			exist = true
			name = f
			return
		}
	}
	err = fmt.Errorf(fileName + ".gitignore not found! Check the list with the -list option!")
	return
}

// constructs an URL with the GitHub raw page
func createUrl(fileName string) (url string) {
	url = strings.Join([]string{GHRawUrl, fileName, GitIgnoreExt}, "")
	return
}

// returns the contents from a web page
func downloadFileContent(url string) (bodyBytes []byte, errReadAll error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, errReadAll = ioutil.ReadAll(resp.Body)
	if errReadAll != nil {
		return
	}
	return
}

// checks if a determined file exists
func fileExists(path string) (exit bool, err error) {
	if _, errStat := os.Stat(path); os.IsNotExist(errStat) {
		exit = true
		return
	}
	err = fmt.Errorf("File already exists")
	return
}

// write some content to a file
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

		_, name, gitIgnoreExistsError := gitIgnoreExists(arg)
		if gitIgnoreExistsError == nil {
			url := createUrl(name)
			content, errDownloadFileContent := downloadFileContent(url)
			if errDownloadFileContent != nil {
				fmt.Println(errDownloadFileContent)
			}
			writeFileContent(content)
		} else {
			fmt.Println(gitIgnoreExistsError)
			os.Exit(1)
		}
	} else {
		flag.Usage()
	}
}
