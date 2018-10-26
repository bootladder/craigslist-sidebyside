package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var urlsFilename = "/home/steve/craigslist-urls"

func saveURL(url string) {
	printf("I'm saving the url " + url)
	touch(urlsFilename)
	contents, err := ioutil.ReadFile(urlsFilename)
	fatal(err)
	newcontents := fmt.Sprintf("%s%s\n", string(contents), url)
	ioutil.WriteFile(urlsFilename, []byte(newcontents), 755)
}

func loadURLs() (urls []string) {
	printf("I'm loading the URLs")

	touch(urlsFilename)
	contents, err := ioutil.ReadFile(urlsFilename)
	fatal(err)
	scanner := bufio.NewScanner(strings.NewReader(string(contents)))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		urls = append(urls, scanner.Text())
	}
	return
}
func touch(filename string) {
	var f, err = os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	fatal(err, "Open "+filename)
	f.Close()
}
