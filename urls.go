package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var urlsFilename = "/home/steve/craigslist-urls"
var urlsFilenameSet2 = "/home/steve/housing.craigslisturls"
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
	log.Print("Saving URLs")
	touch(urlsFilename)
	var newcontents string
	for _, url := range urls {
		newcontents = fmt.Sprintf("%s%s\n", string(newcontents), url)
	}
	ioutil.WriteFile(urlsFilename, []byte(newcontents), 755)
}

func uRLSsetURLAt(index int, url string) {
	urls[index] = url
	save()
}
func uRLSdeleteURLAt(index int) {
	urls = append(urls[:index], urls[(index+1):]...)
	save()
}
func uRLSaddURL() {
	urls = append(urls, "http://denver.craigslist.org/search/sss?query=bike")
	save()
}

func uRLSgetUrls() []string {
	return urls
}

func uRLSloadURLs() {
	log.Print("loadURLs() : ")

	touch(urlsFilename)
	contents, err := ioutil.ReadFile(urlsFilename)
	fatal(err)

	urls = nil //clear the slice

	scanner := bufio.NewScanner(strings.NewReader(string(contents)))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		urls = append(urls, scanner.Text())
	}
}

func uRLSloadURLSet2() {
	log.Print("I'm loading the URLs Set 2")

	touch(urlsFilenameSet2)
	contents, err := ioutil.ReadFile(urlsFilenameSet2)
	fatal(err)

	urls = nil //clear the slice

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
