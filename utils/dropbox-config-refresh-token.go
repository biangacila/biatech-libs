package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type DropboxConfig struct {
	AppKey       string
	AppSecret    string
	RefreshToken string
	AccessToken  string
}

func (config *DropboxConfig) refreshAccessToken() error {
	tokenURL := "https://api.dropboxapi.com/oauth2/token"
	data := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", config.RefreshToken)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(config.AppKey, config.AppSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %s, body: %s", resp.Status, body)
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	config.AccessToken = response.AccessToken
	return nil
}
func getEnvWithDefault(env, defaultValue string) string {
	if os.Getenv(env) != "" {
		return os.Getenv(env)
	}
	return defaultValue
}

func StartRefreshDropboxTokenOnce(clientId, secret, refreshToken string) (newToken string, err error) {
	config := DropboxConfig{
		AppKey:       getEnvWithDefault("DROPBOX_APP_KEY", clientId),
		AppSecret:    getEnvWithDefault("DROPBOX_APP_SECRET", secret),
		RefreshToken: getEnvWithDefault("DROPBOX_REFRESH_TOKEN", refreshToken),
	}
	err = config.refreshAccessToken()
	if err != nil {
		log.Printf("Error refreshing access token: %v \n", err)
		return "", err
	} else {
		newToken = config.AccessToken
	}
	return newToken, nil
}
