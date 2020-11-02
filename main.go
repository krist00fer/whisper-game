package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// To be removed
type Post struct {
	Title  string `json:"Rubrik"`
	Author string `json:"Author"`
	Text   string `json:"Text"`
}

// Whisper describes the messages being received or sent
type Whisper struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

// To be removed
func PostsHandler(w http.ResponseWriter, r *http.Request) {

	posts := []Post{
		Post{"Post one", "Paige", "This is first post."},
		Post{"Post two", "Rachel", "This is second post."},
		Post{"Post three", "Olivia", "This is another post."},
		Post{"Post four", "Kristofer", "This is the last post."},
	}

	json.NewEncoder(w).Encode(posts)
}

func newMessagePosted(w http.ResponseWriter, r *http.Request) {
	log.Println("Message received")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var whisper Whisper
	json.Unmarshal(reqBody, &whisper)
	log.Println("From:", whisper.Sender, "Message:", whisper.Message)

	fmt.Fprintf(w, "Thanks for your message\n")
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/posts", PostsHandler)
	router.HandleFunc("/", newMessagePosted).Methods("POST")

	log.Println("Listening to requests on port 5051")
	log.Fatal(http.ListenAndServe(":5051", router))
}

func main() {
	log.Println("Whisper Service Started")

	handleRequests()
}

/*
	{
		"sender":"olivia",
		"message":"Fun message"
	}
*/
