package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/bitbucket"
	"github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/models"
)

func main() {

	var request models.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	check(err)

	args := os.Args
	inputDir := args[1]

	token, err := bitbucket.RequestToken(request.Source.Key, request.Source.Secret)
	check(err)

	Commit, err := ioutil.ReadFile(inputDir + "/" + request.Params.Commit)
	check(err)

	UpdateCommit := string(Commit)
	UpdateCommit = strings.TrimSpace(UpdateCommit)
	switch state := request.Params.State; state {
	case "success":
		err = bitbucket.SetBuildStatus(request.Source.URL, token, request.Source.APIVersion, request.Source.Team, request.Source.Repo, UpdateCommit, "SUCCESSFUL", request.Source.ConcourseURL)
		check(err)
		log.Print(UpdateCommit)
	case "failed":
		err = bitbucket.SetBuildStatus(request.Source.URL, token, request.Source.APIVersion, request.Source.Team, request.Source.Repo, UpdateCommit, "FAILED", request.Source.ConcourseURL)
		check(err)
		log.Print(UpdateCommit)
	default:
		log.Fatal("No Status Set")
	}

	version := models.MetadataField{Name: "Version", Value: request.Version.Commit}

	metadata := models.Metadata{version}

	err = json.NewEncoder(os.Stdout).Encode(models.OutResponse{Version: request.Version, Metadata: metadata})
	check(err)
}
func check(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
