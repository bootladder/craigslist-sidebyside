package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type messageHandler struct {
	message string
}

func (m *messageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, m.message)
}

func messageHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hellomessagefunc")
}

func createMessageHandlerFunc(msg string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, msg)
	})
}

type Note struct {
	Hello     string    `json:"hello"`
	SearchURL string    `json:"searchURL"`
	CreatedOn time.Time `json:"createdon"`
}

//HTTP Get - /api/notes
func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
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
	// Decode the incoming Note json
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	note.CreatedOn = time.Now()
	//	k := strconv.Itoa(id)
	//	noteStore[k] = note
	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")

	http.ListenAndServe(":8080", r)
}

func main2() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)
	mux.Handle("/welcome", &messageHandler{"hello123"})

	mux.Handle("/func", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "super lambda2")
		}))

	mux.Handle("/closure", createMessageHandlerFunc("my closure thing"))

	http.ListenAndServe(":8080", mux)
}
