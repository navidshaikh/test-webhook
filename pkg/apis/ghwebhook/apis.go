package ghwebhook

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	gh "github.com/navidshaikh/test-webhook/pkg/github"
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
	case *github.IssueCommentEvent:
		if e.Issue.IsPullRequest() && e.Action != nil && *e.Action == "created" && *e.Issue.State == "open" {
			log.Printf("Open PR comment %s \n", *e.Action)
			var ic github.IssueCommentEvent
			if err := json.Unmarshal(payload, &ic); err != nil {
				log.Printf("error loading the webhook payload")
				return
			}
			log.Println(ic.Issue.GetNumber())
			log.Println(ic.Comment.GetBody())
			log.Println(ic.Comment.User.GetLogin())
			log.Println(ic.Comment.GetID())

			ctx := context.Background()
			g, err := gh.DefaultGithub(ctx, "navidshaikh", "test-webhook")
			if err != nil {
				log.Println(err)
			}
			commit, err := g.GetLatestPRCommit(ctx, 6)
			if err != nil {
				log.Println(err)
			}
			log.Println(commit)

		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
