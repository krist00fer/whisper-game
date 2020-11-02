package main

import (
	"encoding/json"
	"net/http"
)

type Post struct {
	Title  string `json:"Rubrik"`
	Author string `json:"Author"`
	Text   string `json:"Text"`
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {

	posts := []Post{
		Post{"Post one", "Paige", "This is first post."},
		Post{"Post two", "Rachel", "This is second post."},
		Post{"Post three", "Olivia", "This is another post."},
		Post{"Post four", "Kristofer", "This is the last post."},
	}

	json.NewEncoder(w).Encode(posts)
}

func main() {
	http.HandleFunc("")
	http.HandleFunc("/posts", PostsHandler)
	http.ListenAndServe(":5051", nil)
}

/*
	{
		"sender":"olivia",
		"message":"Fun message"
	}
*/
