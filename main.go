package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

func main() {
	f, err := os.OpenFile("/tmp/tkg5360.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.Println("tftesting started, listening for webhook..")
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte("addadadda"))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	var data map[string]interface{}
	switch e := event.(type) {
	case *github.IssueCommentEvent:
		log.Printf("Issue comment %s \n", *e.Action)
		if e.Issue.IsPullRequest() && e.Action != nil && *e.Action == "created" && *e.Issue.State == "open" {
			fmt.Println(*e.Comment.ID)
			fmt.Println(*e.Comment.Body)
			fmt.Println(*e.Repo.PullsURL)
			err := json.Unmarshal(payload, &data)
			if err != nil {
				log.Printf("error loading the webhook payload")
				return
			}
			fmt.Println(data["issue"])
			fmt.Println(data["comment"])
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
