package api

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	api "code.gitea.io/gitea/modules/structs"
	"github.com/go-resty/resty/v2"
)

var secret = os.Getenv("SECRET")
var githubToken = "token " + os.Getenv("GH_TOKEN")
var client = resty.New().SetHeader("Authorization", githubToken)

func HookHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	// get the hook event from the headers
	event := r.Header.Get("X-Gitea-Event")

	// only push events are current supported
	if event != "push" {
		log.Printf("received unknown event \"%s\"\n", event)
		return
	}

	// read request body
	var data, err = ioutil.ReadAll(r.Body)
	panicIf(err, "while reading request body")

	// unmarshal request body
	var hook api.PushPayload
	err = json.Unmarshal(data, &hook)
	panicIf(err, fmt.Sprintf("while unmarshaling request base64(%s)", b64.StdEncoding.EncodeToString(data)))

	log.Printf("received webhook on %s", hook.Repo.FullName)

	// find matching config for repository name
	// check if the secret in the configuration matches the request
	if secret != hook.Secret {
		return
	}
	_, err = client.R().
		SetBody(fmt.Sprintf(`{"event_type":"%s push"}`, hook.Repo.FullName)).
		Post("https://api.github.com/repos/Trim21/actions-cron/dispatches")
	if err != nil {
		log.Println("error when dispatch event to github", err.Error())
	}

}

func panicIf(err error, what ...string) {
	if err != nil {
		if len(what) == 0 {
			panic(err)
		}

		panic(errors.New(err.Error() + (" " + what[0])))
	}
}
