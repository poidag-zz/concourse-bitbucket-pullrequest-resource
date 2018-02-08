package bitbucket

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/models"
	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
)

// doObject will request and return a GenericResponse
func doObject(request *http.Request) (*models.GenericResponse, error) {
	var response models.GenericResponse
	err := do(request, &response)
	return &response, err
}

// doSlice will iterate over PagedGenericResponses to return a slice of GenericResponses
func doSlice(request *http.Request) (*[]models.GenericResponse, error) {
	var values []models.GenericResponse
	for {
		var response models.PagedGenericResponse
		err := do(request, &response)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve values")
		}
		values = append(values, response.Values...)
		// fmt.Printf("values: %+v size: %+v page: %+v pagelen: %+v next: %+v\n", len(values), response.Size, response.Page, response.Pagelen, response.Next)
		if response.Next == "" {
			break
		}
		request.URL, err = url.Parse(response.Next)
		if err != nil {
			return &values, errors.Wrapf(err, "failed to parse next url")
		}
	}
	return &values, nil
}

// do will perform a http request with retries and backoff
// will then unmarshall into the passed response object
func do(request *http.Request, response interface{}) error {
	client := pester.New()
	client.MaxRetries = 10
	client.Backoff = pester.ExponentialBackoff
	client.KeepLog = true
	resp, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, client.LogString())
	}

	// Some Bitbucket APIs can return 404 in some cases.
	if resp.StatusCode == 404 {
		return nil
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response body")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.Errorf("request failed, code [%d], url [%s], body: %s", resp.StatusCode, request.URL, buf.String())
	}
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal the response: %s", buf.String())
	}
	return nil
}
