package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var urlsFilename = "/home/steve/craigslist-urls"
var urls []string

func saveURL(url string) {
	printf("I'm saving the url " + url)
	touch(urlsFilename)
	contents, err := ioutil.ReadFile(urlsFilename)
	fatal(err)
	newcontents := fmt.Sprintf("%s%s\n", string(contents), url)
	ioutil.WriteFile(urlsFilename, []byte(newcontents), 755)
}

func save() {
	printf("Saving URLs")
	touch(urlsFilename)
	var newcontents string
	for _, url := range urls {
		newcontents = fmt.Sprintf("%s%s\n", string(newcontents), url)
	}
	ioutil.WriteFile(urlsFilename, []byte(newcontents), 755)
}

func setURLAt(index int, url string) {
	urls[index] = url
	save()
}

func getUrls() []string {
	return urls
}

func loadURLs() {
	printf("I'm loading the URLs")

	touch(urlsFilename)
	contents, err := ioutil.ReadFile(urlsFilename)
	fatal(err)
	scanner := bufio.NewScanner(strings.NewReader(string(contents)))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		urls = append(urls, scanner.Text())
	}
}
func touch(filename string) {
	var f, err = os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	fatal(err, "Open "+filename)
	f.Close()
}
