package gitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	githubV3BaseUrl      = "https://api.github.com/"
	githubV3AcceptHeader = "application/vnd.github.v3+json"
)

type Release struct {
	Url        string `json:"url"`
	HtmlUrl    string `json:"html_url"`
	Name       string `json:"name"`
	TarballUrl string `json:"tarball_url"`
	ZipballUrl string `json:"zipball_url"`
	Body       string `json:"body"`
}

func ParseBody(bodyBlob []byte) (Release, error) {
	var latestRelease Release
	err := json.Unmarshal(bodyBlob, &latestRelease)
	return latestRelease, err
}

func GetLatestReleaseUrl(owner string, repo string) (string, error) {
	latestRelease, err := GetLatestRelease(owner, repo)
	if err != nil {
		return "", err
	}
	latestUrl := latestRelease.TarballUrl
	if latestUrl == "" {
		message := fmt.Sprintf("No tarball found for %s", repo)
		return "", errors.New(message)
	}
	return latestUrl, nil
}

func GetLatestRelease(owner string, repo string) (Release, error) {
	var emptyRelease Release

	// GET https://api.github.com/repos/:owner/:repo/releases/latest
	latestReleaseUrl := fmt.Sprintf("%srepos/%s/%s/releases/latest",
		githubV3BaseUrl, owner, repo)

	client := &http.Client{}
	req, err := http.NewRequest("GET", latestReleaseUrl, nil)
	req.Header.Add("Accept", githubV3AcceptHeader)
	response, err := client.Do(req)
	if err != nil {
		return emptyRelease, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return emptyRelease, errors.New(response.Status)
	}
	bodyBlob, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return emptyRelease, err
	}

	latest, err := ParseBody(bodyBlob)
	return latest, err
}
