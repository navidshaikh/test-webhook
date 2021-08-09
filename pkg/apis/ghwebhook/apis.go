package ghwebhook

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/webhook", handleWebhook)
	log.Println("tftesting started, listening for webhook..")
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

	switch e := event.(type) {
	/*
		case *github.PullRequestReviewCommentEvent:
			log.Printf("Pull request review comment\n", *e.Action)
		case *github.PullRequestReviewEvent:
			log.Prinf("Pull request review\n", *e.Action)
	*/
	case *github.IssueCommentEvent:
		if e.Issue.IsPullRequest() && e.Action != nil && *e.Action == "created" && *e.Issue.State == "open" {
			log.Printf("Issue comment %s \n", *e.Action)
			var ic github.IssueComment
			if err := json.Unmarshal(payload, &ic); err != nil {
				log.Printf("error loading the webhook payload")
				return
			}

		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
