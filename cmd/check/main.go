package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/bitbucket"
	"github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/models"
)

func main() {
	var response models.CheckResponse
	var request models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	check(err)

	token, err := bitbucket.RequestToken(request.Source.Key, request.Source.Secret)
	check(err)

	out, err := bitbucket.GetPullRequests(request.Source.URL, token, request.Source.APIVersion, request.Source.Team, request.Source.Repo)
	check(err)

	counter := 0
	for counter < 1 {
		for _, pr := range *out {

			state, err := bitbucket.GetCommitStatus(pr.Source.Commit.Links.Self.Href, token)
			check(err)

			link := pr.Links.HTML.Href

			if pr.CommentCount > 0 {
				comments, err := bitbucket.GetPrComments(pr.Links.Comments.Href, token)
				check(err)

				for _, comment := range comments {

					possibleCommand := strings.Split(comment.Content.Raw, "\n")[0]

					// If the first line of the comment is "/retest", then include this link
					// in the output, instead of the default PR link. This should trigger
					// a new build.
					if possibleCommand == "/retest" {
						link = comment.Link
					}
				}
			}

			responseOut := models.Version{
				Commit:      pr.Source.Commit.Hash,
				PullRequest: strconv.Itoa(pr.ID),
				Link:        link,
			}

			switch state {
			case "SUCCESSFUL":
				response = append(response, responseOut)
			case "INPROGRESS":
				response = append(response, responseOut)
			case "FAILING", "FAILED":
				response = append(response, responseOut)
				counter++
			case "STOPPED":
				response = append(response, responseOut)
				counter++
			case "none":
				response = append(response, responseOut)
				counter++
			default:
				counter++
			}
		}
		break
	}
	err = json.NewEncoder(os.Stdout).Encode(response)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
