package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/models"
)

// SetBuildStatus updates the commit associated with a pull-request and sets the state () as well as a link to the Concourse build log.
func SetBuildStatus(url, token, version, team, repo, commit, state, concourseHost string) error {
	if url == "" {
		return errors.New("url must be provided")
	}
	if token == "" {
		return errors.New("token must be provided")
	}
	if version == "" {
		return errors.New("version must be provided")
	}
	if team == "" {
		return errors.New("team must be provided")
	}
	if repo == "" {
		return errors.New("repo must be provided")
	}
	if commit == "" {
		return errors.New("commit must be provided")
	}
	if state == "" {
		return errors.New("state must be provided")
	}
	if concourseHost == "" {
		return errors.New("concourse host must be provided")
	}

	buildJob := os.Getenv("BUILD_JOB_NAME")

	concourseURL := fmt.Sprintf(
		"%s/teams/%s/pipelines/%s/jobs/%s/builds/%s",
		concourseHost,
		os.Getenv("BUILD_TEAM_NAME"),
		os.Getenv("BUILD_PIPELINE_NAME"),
		buildJob,
		os.Getenv("BUILD_NAME"),
	)

	key := "concourse-" + buildJob

	status := models.OutStatus{State: state, Key: key, URL: concourseURL}
	out, err := json.Marshal(status)
	if err != nil {
		return errors.Wrapf(err, "unable to marshal build status: %+v", status)
	}

	req, err := http.NewRequest("POST", url+"/"+version+"/repositories/"+team+"/"+repo+"/commit/"+commit+"/statuses/build", bytes.NewBuffer(out))
	if err != nil {
		return errors.Wrap(err, "unable to create request object")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "request to set build status failed")
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			return errors.Wrapf(err, "request to set build status failed with status [%d], but the response body could not be read", res.StatusCode)
		}
		return errors.Errorf("request to set build status failed, code [%d], url [%s], body: %s", res.StatusCode, req.URL, buf.String())
	}
	return nil
}

// GetPullRequests fetches the pull requests for a specific repository.
func GetPullRequests(url string, token string, version string, team string, repo string) (*[]models.GenericResponse, error) {
	if url == "" {
		return nil, errors.New("url must be provided")
	}
	if token == "" {
		return nil, errors.New("token must be provided")
	}
	if version == "" {
		return nil, errors.New("version must be provided")
	}
	if team == "" {
		return nil, errors.New("team must be provided")
	}
	if repo == "" {
		return nil, errors.New("repo must be provided")
	}

	req, err := http.NewRequest("GET", url+"/"+version+"/repositories/"+team+"/"+repo+"/pullrequests", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request object")
	}
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := doSlice(req)
	if err != nil {
		return nil, errors.Wrap(err, "request to retrieve pull requests failed")
	}
	return response, nil
}

// GetCommitStatus retrieves the current commit status for a specific commit, referenced by URL.
func GetCommitStatus(url string, token string) (string, error) {
	// Ref <https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/pullrequests/%7Bpull_request_id%7D/statuses>

	if url == "" {
		return "", errors.New("url must be provided")
	}
	if token == "" {
		return "", errors.New("token must be provided")
	}

	req, err := http.NewRequest("GET", url+"/statuses", nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to create request")
	}
	req.Header.Add("Authorization", "Bearer "+token)

	var response models.CommitResponse
	err = do(req, &response)
	if err != nil {
		return "", errors.Wrap(err, "request to retrieve commit status failed")
	}
	if len(response.Values) > 0 {
		return response.Values[0].State, nil
	}
	return "none", nil
}

// GetPrComments returns the comments associated with a specific pullrequest, referenced by URL.
func GetPrComments(url string, token string) (comments []models.Comment, err error) {
	// Ref https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/pullrequests/%7Bpull_request_id%7D/comments

	if url == "" {
		return comments, errors.New("url must be provided")
	}
	if token == "" {
		return comments, errors.New("token must be provided")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return comments, errors.Wrap(err, "unable to create request")
	}
	req.Header.Add("Authorization", "Bearer "+token)

	var response models.PagedGenericResponse
	err = do(req, &response)
	if err != nil {
		return comments, errors.Wrap(err, "request to retrieve commit status failed")
	}

	if len(response.Values) > 0 {

		for _, commentRef := range response.Values {

			if commentRef.Inline != nil {
				// skip over inline comments
				continue
			}

			if commentRef.Parent != nil {
				// If its a reply to another comment, ignore it too.
				continue
			}

			comments = append(comments, models.Comment{
				User:      commentRef.User,
				Content:   commentRef.Content,
				CreatedOn: commentRef.CreatedOn,
				Link:      commentRef.Links.HTML.Href,
			})

		}

	}
	return comments, nil
}

func GetPullRequestByID(url string, token string, version string, team string, repo string, request string) (*models.GenericResponse, error) {
	if url == "" {
		return nil, errors.New("url must be provided")
	}
	if token == "" {
		return nil, errors.New("token must be provided")
	}
	if version == "" {
		return nil, errors.New("version must be provided")
	}
	if team == "" {
		return nil, errors.New("team must be provided")
	}
	if repo == "" {
		return nil, errors.New("repo must be provided")
	}
	if request == "" {
		return nil, errors.New("PR id must be provided")
	}

	req, err := http.NewRequest("GET", url+"/"+version+"/repositories/"+team+"/"+repo+"/pullrequests/"+request, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request object")
	}
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := doObject(req)
	if err != nil {
		return nil, errors.Wrap(err, "request to retrieve pull request failed")
	}
	return response, nil
}

func ApprovePullRequest(url string, token string, version string, team string, repo string, request string) (*models.GenericResponse, error) {
	if url == "" {
		return nil, errors.New("url must be provided")
	}
	if token == "" {
		return nil, errors.New("token must be provided")
	}
	if version == "" {
		return nil, errors.New("version must be provided")
	}
	if team == "" {
		return nil, errors.New("team must be provided")
	}
	if repo == "" {
		return nil, errors.New("repo must be provided")
	}
	if request == "" {
		return nil, errors.New("PR id must be provided")
	}

	req, err := http.NewRequest("POST", url+"/"+version+"/repositories/"+team+"/"+repo+"/pullrequests/"+request+"/approve", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request object")
	}
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := doObject(req)
	if err != nil {
		return nil, errors.Wrap(err, "request to approve pull request failed")
	}
	return response, nil
}

func DeclinePullRequest(url string, token string, version string, team string, repo string, request string) (*models.GenericResponse, error) {
	if url == "" {
		return nil, errors.New("url must be provided")
	}
	if token == "" {
		return nil, errors.New("token must be provided")
	}
	if version == "" {
		return nil, errors.New("version must be provided")
	}
	if team == "" {
		return nil, errors.New("team must be provided")
	}
	if repo == "" {
		return nil, errors.New("repo must be provided")
	}
	if request == "" {
		return nil, errors.New("PR id must be provided")
	}

	req, err := http.NewRequest("POST", url+"/"+version+"/repositories/"+team+"/"+repo+"/pullrequests/"+request+"/decline", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request object")
	}
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := doObject(req)
	if err != nil {
		return nil, errors.Wrap(err, "request to decline pull request failed")
	}

	return response, nil
}

func RequestToken(key string, secret string) (string, error) {
	if key == "" {
		return "", errors.New("key must be provided")
	}
	if secret == "" {
		return "", errors.New("secret must be provided")
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", "https://bitbucket.org/site/oauth2/access_token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "unable to create request object")
	}
	req.SetBasicAuth(key, secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var response models.Token
	err = do(req, &response)
	if err != nil {
		return "", errors.Wrap(err, "request for token failed")
	}
	return response.AccessToken, nil
}
