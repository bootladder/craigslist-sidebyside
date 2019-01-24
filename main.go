package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

var debug = true

var err error

var urlstore urlStore

type craigslistPostRequest struct {
	SearchURL   string `json:"searchURL"`
	ColumnIndex int    `json:"columnIndex"`
	SetIndex    int    `json:"setIndex"`
}
type craigslistPostResponse struct {
	ResponseHTML string   `json:"response"`
	Urls         []string `json:"urls"`
}

type craigslistDeleteRequest struct {
	ColumnIndex int `json:"columnIndex"`
	SetIndex    int `json:"setIndex"`
}
type craigslistAddRequest struct {
	SetIndex int `json:"setIndex"`
}
type craigslistGetRequest struct {
	SetIndex int `json:"setIndex"`
}
type returnURLSetResponse struct {
	Urls []string `json:"urls"`
}

func main() {

	inject()
	urlstore.loadURLs()

	router := httprouter.New()
	router.ServeFiles("/static/*filepath",
		http.Dir("public"))
	router.ServeFiles("/images/*filepath",
		http.Dir("public/images"))

	router.POST("/api/", createPostHandler(""))
	router.GET("/api/:setIndex", getURLSet)
	router.GET("/api/", getAllURLSetNames)
	router.DELETE("/api/", createDeleteHandler(""))
	router.PUT("/api/", createPutHandler(""))

	browser.OpenURL("http://localhost:8080/static/index.html")
	http.ListenAndServe(":8080", router)
}

func postNoteHandler(w http.ResponseWriter, r *http.Request) {

	var req craigslistPostRequest = parsePostRequestBody(r.Body)

	var resp craigslistPostResponse
	resp.ResponseHTML = fetchCraigslistQuery(req.SearchURL)
	urlstore.setURLAt(req.SetIndex, req.ColumnIndex, req.SearchURL)
	resp.Urls = urlstore.urlsets[req.SetIndex]

	writePostResponse(w, resp)
}

func parsePostRequestBody(requestBody io.Reader) craigslistPostRequest {
	var req craigslistPostRequest
	err := json.NewDecoder(requestBody).Decode(&req)
	fatal(err)

	req.SearchURL, err = url.QueryUnescape(req.SearchURL)
	fatal(err)

	return req
}

func fetchCraigslistQuery(url string) string {
	if debug == true {
		return `<html><body><ul><li class="result-row" data-pid="6744258112">` +
			` Wow cool ` + url + ` </li></ul></body></html>`

	}
	rawHtml, err := makeRequest(url)
	if err != nil {
		return `<html><body><ul><li class="result-row" data-pid="6744258112">` +
			` ERROR: ` + err.Error() + ` : ` + url + ` </li></ul></body></html>`
	}

	return extractCraigslistResultRows(rawHtml)
}

func extractCraigslistResultRows(rawHtml string) string {

	doc, _ := html.Parse(strings.NewReader(rawHtml))
	resultRows, _ := getResultRows(doc)
	return renderNode(resultRows)
}

func getResultRows(doc *html.Node) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "result-row" {
					b = n.Parent
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Missing <result rows> in the node tree")
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func writePostResponse(w http.ResponseWriter,
	resp craigslistPostResponse) {
	jsonOut, err := json.Marshal(resp)
	fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOut)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	var req craigslistDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	fatal(err)
	fmt.Printf("Delete: the index is %d\n", req.ColumnIndex)

	urlstore.deleteURLAt(req.SetIndex, req.ColumnIndex)

	returnURLSetJSONResponse(w, req.SetIndex)
}

func putHandler(w http.ResponseWriter, r *http.Request) {

	var req craigslistAddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	fatal(err, "JSON Decode Put Handler Body")
	fmt.Printf("Add(Put): the index is %d\n", req.SetIndex)
	urlstore.addURL(req.SetIndex)
	returnURLSetJSONResponse(w, req.SetIndex)
}

func getURLSet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	setIndexString := ps.ByName("setIndex")
	setIndex, err := strconv.Atoi(setIndexString)
	fatal(err)

	if len(urlstore.urlsets) < (setIndex + 1) {
		urlstore.addNewURLSet()
	}

	returnURLSetJSONResponse(w, setIndex)
}

func getAllURLSetNames(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var names []string = urlstore.getAllURLSetNames()

	jsonOut, err := json.MarshalIndent(names, "", "  ")
	fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOut)
}

func makeRequest(url string) (string, error) {

	log.Print("makeRequest: sleep ... ")
	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)

	log.Printf("makeRequest: %s\n", url)

	client := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Get(url) //"https://httpbin.org/get"

	//gracefully handle error with invalid craigslist URL
	if err != nil {
		log.Println("    TIMEOUT: " + url)
		return "TIMEOUT", errors.New("TIMEOUT")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body), nil
}

func returnURLSetJSONResponse(w http.ResponseWriter, setIndex int) {
	var resp returnURLSetResponse
	resp.Urls = urlstore.urlsets[setIndex]

	jsonOut, err := json.MarshalIndent(resp, "", "  ")
	fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOut)
}

func fatal(err error, msgs ...string) {
	if err != nil {
		var str string
		for _, msg := range msgs {
			str = msg
			break
		}
		panic(errors.Wrap(err, str))
	}
}

func printf(s string, a ...interface{}) {
	fmt.Printf(s, a...)
}

func createPostHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		postNoteHandler(w, r)
	}
}

func createDeleteHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		deleteHandler(w, r)
	}
}

func createPutHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		putHandler(w, r)
	}
}

var external = externalFuncs{}

type externalFuncs struct {
	readfile  func(string) ([]byte, error)
	writefile func(string, []byte, os.FileMode) error
}

func inject() {
	external.readfile = ioutil.ReadFile
	external.writefile = ioutil.WriteFile
}
