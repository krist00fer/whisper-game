package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const version = "0.0.6"

var sender string
var forwardAddress string
var gossip []Whisper

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

	gossip = []Whisper{}

	log.Println("Whisper Service (version:", version, ") - Started")
	log.Println("Sender:", sender)
	log.Println("Forward Address:", forwardAddress)

	setupAndHandleRequests()
}

func setupAndHandleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/config", handleGetConfig)
	router.HandleFunc("/whisper", handlePostMessage).Methods("POST")
	router.HandleFunc("/gossip", handleGetGossip)

	log.Println("Listening to requests on port 5051")
	log.Fatal(http.ListenAndServe(":5051", router))
}

func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("Message Received")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var whisper Whisper
	json.Unmarshal(reqBody, &whisper)
	log.Println("From:", whisper.Sender, "Message:", whisper.Message)

	gossip = append(gossip, whisper)

	if whisper.Sender == sender {
		log.Println("Whisper came back")

	} else {
		sendWhisper(whisper)
	}
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

func handleGetGossip(w http.ResponseWriter, r *http.Request) {
	log.Println("Gossig Requested")

	json.NewEncoder(w).Encode(gossip)
}

func sendWhisper(w Whisper) {
	log.Println("Forwarding whisper to", forwardAddress)

	jsonReq, _ := json.Marshal(w)
	resp, err := http.Post(forwardAddress, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln("Unable to forward whisper to", forwardAddress)
		log.Fatalln(err)
	}

	defer resp.Body.Close()
}

// Whisper describes the messages being received or sent
type Whisper struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

/*
	{
		"sender":"olivia",
		"message":"Fun message"
	}
*/
