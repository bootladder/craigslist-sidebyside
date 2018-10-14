package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

var err error

type Note struct {
	Hello     string    `json:"hello"`
	SearchURL string    `json:"searchURL"`
	CreatedOn time.Time `json:"createdon"`
}

func createGetHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		GetNoteHandler(w, r)
	}
}
func createPostHandler(msg string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		PostNoteHandler(w, r)
	}
}

//HTTP Get - /api/notes
func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("hello\n")
	var messages []string
	messages = append(messages, "arr1")
	messages = append(messages, "arr2")
	messages = append(messages, "arr3")
	messages = append(messages, "arr4")

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(messages)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//HTTP Post - /api/notes
func PostNoteHandler(w http.ResponseWriter, r *http.Request) {

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	note.CreatedOn = time.Now()
	note.SearchURL, _ = url.QueryUnescape(note.SearchURL)

	j, err := json.Marshal(note)
	Fatal(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func main() {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath",
		http.Dir("/home/steve/prog/go/src/scratch/stdlib-http-server/public"))

	router.POST("/api/", createPostHandler(""))

	browser.OpenURL("http://localhost:8080/static/index.html")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Fatal panics on error
//First parameter of msgs is used each following variadic arg is dropped
func Fatal(err error, msgs ...string) {
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
	fmt.Printf(s, a)
}
