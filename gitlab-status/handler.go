package function

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/openfaas/openfaas-cloud/sdk"
)

// Handle a serverless request
func Handle(req []byte) string {
	status, statusErr := sdk.UnmarshalStatus(req)
	if statusErr != nil {
		return statusErr.Error()
	}

	token, tokenErr := sdk.ReadSecret("gitlab-api-token")
	if tokenErr != nil {
		return tokenErr.Error()
	}

	url := gitLabURLBuilder(status)

	for _, commitStatus := range status.CommitStatuses {
		reportErr := sendReport(url, token, commitStatus.Status, commitStatus.Description, commitStatus.Context)
		if reportErr != nil {
			log.Fatalf("failed to report status %v, error: %s", status, reportErr.Error())
		}
	}

	return ""
}

func gitLabURLBuilder(status *sdk.Status) string {
	delimeterPlace := 3
	delimeterValue := "/"
	baseURL := getURLbyDelimeter(status.EventInfo.URL, delimeterPlace, delimeterValue)
	routeURL := fmt.Sprintf("/api/v4/projects/%d/statuses/%s?", status.EventInfo.InstallationID, status.EventInfo.SHA)
	newURL := append(strings.Split(baseURL, ""), routeURL)
	wholeURL := strings.Join(newURL, "")
	return wholeURL
}

func getURLbyDelimeter(baseURL string, position int, delimeter string) string {
	sliceURL := strings.Split(baseURL, "")
	var slash, last int
	for index, symbol := range sliceURL {
		if slash == position {
			last = index
			break
		}
		if symbol == delimeter {
			slash++
		}
	}

	sliceURL = append(sliceURL[:last-1], sliceURL[cap(sliceURL):]...)
	baseURL = strings.Join(sliceURL, "")
	return baseURL
}

func appendParameters(URL string, state string, desc string, context string) (string, error) {
	var theURL *url.URL

	theURL, urlErr := url.Parse(URL)
	if urlErr != nil {
		return "", urlErr
	}

	if state == "failure" {
		state = "failed"
	}

	parameters := url.Values{}
	parameters.Add("state", state)
	parameters.Add("description", desc)
	parameters.Add("context", context)
	theURL.RawQuery = parameters.Encode()

	return theURL.String(), nil

}

func sendReport(URL string, token string, state string, desc string, context string) error {
	fullURL, fullURLErr := appendParameters(URL, state, desc, context)
	if fullURLErr != nil {
		return fullURLErr
	}
	var b io.Reader
	req, reqErr := http.NewRequest("POST", fullURL, b)
	if reqErr != nil {
		return reqErr
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	client := &http.Client{}
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return clientErr
	}
	resp.Body.Close()

	return nil
}
