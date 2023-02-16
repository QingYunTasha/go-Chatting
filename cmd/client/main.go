package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handleChat)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGetChat(w, r)
	} else if r.Method == "POST" {
		handlePostChat(w, r)
	} else {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func handleGetChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<html><body>")
	fmt.Fprintln(w, "<h1>Welcome to the chat room!</h1>")
	fmt.Fprintln(w, "<form method=\"post\">")
	fmt.Fprintln(w, "<input type=\"text\" name=\"message\" size=\"50\">")
	fmt.Fprintln(w, "<input type=\"submit\" value=\"Send\">")
	fmt.Fprintln(w, "</form>")
	fmt.Fprintln(w, "</body></html>")
}

func handlePostChat(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "Empty message", http.StatusBadRequest)
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s: %s\n", timestamp, message)
	fmt.Fprintf(w, "Message sent: %s\n", message)
}
