package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/refto/server/config"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func GetAccessToken(code string) (token string, err error) {
	hc := http.Client{Timeout: 30 * time.Second}
	reqData := []byte(`{"client_id":"` + config.Get().GitHub.ClientID + `", "client_secret":"` + config.Get().GitHub.ClientSecret + `", "code":"` + code + `"}`)
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token/", bytes.NewBuffer(reqData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(http.StatusText(resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	data := GetAccessTokenResp{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	if data.Error != "" {
		err = errors.New(data.ErrorDescription)
		return
	}

	if data.AccessToken == "" {
		err = errors.New("access token missing")
		return
	}

	return data.AccessToken, nil
}

func GetUser(token string) (user *github.User, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(http.StatusText(resp.StatusCode))
		return
	}

	return
}
