package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

const (
	redirectURL = "http://localhost:8080" // TODO: dynamically get an available port
)

var (
	dataPath        string
	tokenFile       string
	credentialsFile string
	Scopes          = []string{calendar.CalendarReadonlyScope}
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	dataPath = filepath.Join(homeDir, ".config/cal-term/")
	tokenFile = filepath.Join(dataPath, "token.json")
	credentialsFile = filepath.Join(dataPath, "credentials.json")

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type Credentials struct {
	Installed *CredentialDetails `json:"installed"`
}

type CredentialDetails struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
}

type Auth struct {
	CredentialsFile string
	TokenFile       string
	OAuthConfig     *oauth2.Config
}

func New(clientId, clientSecret string) *Auth {
	return &Auth{
		CredentialsFile: credentialsFile,
		TokenFile:       tokenFile,
		OAuthConfig: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  redirectURL,
			Scopes: []string{
				calendar.CalendarReadonlyScope,
			},
		},
	}
}

func (a *Auth) StoreCredentials() error {
	c := &Credentials{
		Installed: &CredentialDetails{
			ClientID:     a.OAuthConfig.ClientID,
			ClientSecret: a.OAuthConfig.ClientSecret,
			RedirectURIs: []string{redirectURL},
			AuthURI:      google.Endpoint.AuthURL,
			TokenURI:     google.Endpoint.TokenURL,
		},
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(a.CredentialsFile, data, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: program needs sudo/root privileges to write to %s", credentialsFile)
		}
		return err
	}
	return nil
}

func (a *Auth) StoreToken(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = os.WriteFile(a.TokenFile, data, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: program needs sudo/root privileges to write to %s", tokenFile)
		}
		return err
	}
	return nil
}

func (a *Auth) GetTokenFromWeb() (*oauth2.Token, error) {
	ch := make(chan *oauth2.Token)
	done := make(chan bool)

	state := generateState()
	authURL := a.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	server := redirectServer(8080, state, a.OAuthConfig, ch)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v", err)
		}
		done <- true
	}()

	waitForKeyPress(fmt.Sprintf("Press any key to open the auth URL:\n👉 %s", authURL))
	openUrl(authURL)

	select {
	case token := <-ch:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
		<-done
		return token, nil
	case <-time.After(5 * time.Minute):
		server.Close()
		<-done
		return nil, fmt.Errorf("timeout waiting for authentication")
	}
}

func GetTokenFromFile() (*oauth2.Token, error) {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{}
	err = json.Unmarshal(data, token)
	return token, err
}

func GetConfigFromFile() (*oauth2.Config, error) {
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(data, Scopes...)
	return config, err
}

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
