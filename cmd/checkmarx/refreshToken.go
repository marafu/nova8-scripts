package checkmarx

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/marafu/nova8-scripts/cmd/utils"
)

type InputRefreshToken struct {
	ClientID     string
	RefreshToken string
}

type ResponseRefreshToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

func GetRefreshToken(config utils.Config) (*ResponseRefreshToken, error) {

	uri := fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/token", config.Checkmarx.TenantName, config.Checkmarx.AuthUrl)

	payload := strings.NewReader(fmt.Sprintf("grant_type=refresh_token&client_id=%s&refresh_token=%s", config.Checkmarx.ClientID, config.Checkmarx.AccessToken))

	req, err := http.NewRequest("POST", uri, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	proxyUrl, err := url.Parse(config.General.Proxy)

	if err != nil {
		return nil, err
	}

	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		var response ResponseRefreshToken

		err = json.Unmarshal(body, &response)

		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil
}
