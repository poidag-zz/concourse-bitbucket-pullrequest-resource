package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/bitbucket"
	"github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/models"
)

func main() {

	var request models.InRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	check(err)

	if request.Version.Commit == "" || request.Version.PullRequest == "" {
		log.Printf("Ignoring input request without version (commit/PR)")
		err = json.NewEncoder(os.Stdout).Encode(models.InResponse{Version: models.Version{}, Metadata: models.Metadata{}})
		check(err)
		return
	}

	token, err := bitbucket.RequestToken(request.Source.Key, request.Source.Secret)
	check(err)

	out, err := bitbucket.GetPullRequestByID(request.Source.URL, token, request.Source.APIVersion, request.Source.Team, request.Source.Repo, request.Version.PullRequest)
	check(err)

	err = bitbucket.SetBuildStatus(request.Source.URL, token, request.Source.APIVersion, request.Source.Team, request.Source.Repo, out.Source.Commit.Hash, "INPROGRESS", request.Source.ConcourseURL)
	check(err)

	inVersion := request.Version

	args := os.Args

	outputDir := args[1]
	versionID := []byte(request.Version.PullRequest)
	commitID := []byte(string(strings.Replace(out.Source.Commit.Hash, "\n", "", -1)))
	Branch := []byte(string(strings.Replace(out.Source.Branch.Name, "\n", "", -1)))

	err = os.MkdirAll(outputDir, os.ModePerm)
	check(err)

	r, err := git.PlainClone(outputDir, false, &git.CloneOptions{
		URL: "https://x-token-auth:" + token + "@bitbucket.org/" + request.Source.Team + "/" + request.Source.Repo,
	})
	check(err)

	w, err := r.Worktree()
	err = w.Checkout(&git.CheckoutOptions{

		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", out.Source.Branch.Name)),
		Force:  true,
	})

	err = ioutil.WriteFile(outputDir+"/version", versionID, 0644)
	check(err)

	err = ioutil.WriteFile(outputDir+"/commit", commitID, 0644)
	check(err)

	err = ioutil.WriteFile(outputDir+"/branch", Branch, 0644)
	check(err)

	version := models.MetadataField{Name: "Version", Value: request.Version.Commit}
	author := models.MetadataField{Name: "Author", Value: out.Author.DisplayName}
	branch := models.MetadataField{Name: "Branch", Value: out.Source.Branch.Name}
	commit := models.MetadataField{Name: "Commit", Value: out.Source.Commit.Hash}
	metadata := models.Metadata{version, author, branch, commit}

	err = json.NewEncoder(os.Stdout).Encode(models.InResponse{Version: inVersion, Metadata: metadata})
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
