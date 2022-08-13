package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseResource   = "/admin/index.php"
	loginResource  = "route=common/login"
	updateResource = "route=localisation/currency/refresh"
)

// requestCredentials is a container for security data which user for update requests
type requestCredentials struct {
	userToken string
	cookies   []*http.Cookie
}

var client = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
	Timeout: time.Second * 10,
}

func getUserToken() (*requestCredentials, error) {
	form := url.Values{}
	form.Set("username", opts.Login)
	form.Set("password", opts.Password)

	req, err := http.NewRequest("POST", fmt.Sprintf(`%s%s?%s`, opts.BaseOpenCartUrl, baseResource, loginResource), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare new reqeust: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	var body []byte
	if resp != nil {
		body, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}

	// only redirection status code returned when request has required data
	if resp.StatusCode != 302 {
		fmt.Printf("%s", body)
		return nil, fmt.Errorf("site return error")
	}

	// extract user_token from redirect url
	redirectLocation := resp.Header.Get("Location")
	urlLocation, err := url.Parse(redirectLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	userToken := urlLocation.Query()["user_token"][0]
	cookie := resp.Cookies()

	return &requestCredentials{
		userToken: userToken,
		cookies:   cookie,
	}, nil
}

func doUpdateCurrency(credentials *requestCredentials) error {
	updateURL := fmt.Sprintf(`%s%s?%s&user_token=%s`, opts.BaseOpenCartUrl, baseResource, updateResource, credentials.userToken)

	req, err := http.NewRequest("GET", updateURL, nil)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}

	if len(credentials.cookies) > 0 {
		req.AddCookie(credentials.cookies[0])
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	var body []byte
	if resp != nil {
		_, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}

	// only redirection status code 302 returned when request has required data
	if resp.StatusCode != 302 {
		fmt.Printf("%s", body)
		return fmt.Errorf("site return error")
	}

	return nil
}
