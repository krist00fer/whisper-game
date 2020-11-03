package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const version = "0.0.5"

var sender string
var forwardAddress string

func main() {
	flag.StringVar(&sender, "sender", "N/A", "Set the current sender")
	flag.StringVar(&forwardAddress, "forwardAddress", "N/A", "Address to foward whispers to")
	flag.Parse()

	val, present := os.LookupEnv("WHISPER_SENDER")
	if sender == "N/A" && present {
		sender = val
	}

	val, present = os.LookupEnv("WHISPER_FORWARD_ADDRESS")
	if forwardAddress == "N/A" && present {
		forwardAddress = val
	}

	log.Println("Whisper Service (version:", version, ") - Started")
	log.Println("Sender:", sender)
	log.Println("Forward Address:", forwardAddress)

	setupAndHandleRequests()
}

func setupAndHandleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/posts", PostsHandler)
	router.HandleFunc("/config", handleGetConfig)
	router.HandleFunc("/", handlePostMessage).Methods("POST")

	log.Println("Listening to requests on port 5051")
	log.Fatal(http.ListenAndServe(":5051", router))
}

func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("Message Received")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var whisper Whisper
	json.Unmarshal(reqBody, &whisper)
	log.Println("From:", whisper.Sender, "Message:", whisper.Message)

	fmt.Fprintf(w, "Thanks for your message\n")
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Config Requested")
	fmt.Fprintln(w, "Whisper Service - Configuration")
	fmt.Fprintln(w, "----------------------------------------------------------------------------------------------------")
	fmt.Fprintln(w, "Version         :", version)
	fmt.Fprintln(w, "Sender          :", sender)
	fmt.Fprintln(w, "Forward Address :", forwardAddress)
}

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

/*
	{
		"sender":"olivia",
		"message":"Fun message"
	}
*/
