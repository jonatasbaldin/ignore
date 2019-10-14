package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

const (
	GitIgnoreExt    = ".gitignore"
	GHRawUrl        = "https://raw.githubusercontent.com/dvcs/gitignore/master/templates/"
	GHGitignoreRepo = "https://github.com/dvcs/gitignore/tree/master/templates"
)

// return a slice of strings containing all the .gitignore file names from GitHub
func listFiles(pageUrl string) (files []string) {
	c := colly.NewCollector(
		colly.AllowedDomains("github.com", "127.0.0.1"),
	)
	c.OnHTML("table.files td.content a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasSuffix(href, GitIgnoreExt) {
			unscapedHref, _ := url.QueryUnescape(href)
			files = append(files, path.Base(strings.TrimSuffix(unscapedHref, GitIgnoreExt)))
		}
	})

	if err := c.Visit(pageUrl); err != nil {
		log.Fatalf("Couldn't visit the URL %q – [Error]: %q", pageUrl, err)
	}
	return
}

// checks if the specified file name exists in the .gitignore list of files from GitHub
func gitIgnoreExists(listUrl string, fileName string) (exist bool, name string, err error) {
	for _, f := range listFiles(listUrl) {
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
		fmt.Println(strings.Join(listFiles(GHGitignoreRepo), "\n"))
		return
	}

	if len(os.Args) > 1 {
		_, fileExistsError := fileExists(GitIgnoreExt)
		if fileExistsError != nil {
			log.Fatal(fileExistsError)
		}

		var content []byte
		for i, arg := range os.Args[1:] {
			_, name, err := gitIgnoreExists(GHGitignoreRepo, arg)
			if err != nil {
				log.Fatal(err)
			}

			url := createUrl(name)
			argContent, err := downloadFileContent(url)
			if err != nil {
				log.Fatal(err)
			}

			if i != 0 {
				// add a blank line for every arg except the first
				content = append(content, byte('\n'))
			}

			// add a header with the url of the arg
			header := fmt.Sprintf("# %s\n\n", url)
			content = append(content, []byte(header)...)

			content = append(content, argContent...)
		}

		writeFileContent(content)
	} else {
		flag.Usage()
	}
}
