package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/github"
    "io/ioutil"
)

var (
    githubOAuthConfig = &oauth2.Config{
        ClientID:     "",     // Replace this
        ClientSecret: "", // Replace this
        RedirectURL:  "https://expert-guacamole-9j5ppg5j7qrc7v5p-8080.app.github.dev/callback",
        Scopes:       []string{"user:email"},
        Endpoint:     github.Endpoint,
    }

    oauthStateString = "random-state" // CSRF protection
)

func main() {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/login", handleGitHubLogin)
    http.HandleFunc("/callback", handleGitHubCallback)

    fmt.Println("App running at https://expert-guacamole-9j5ppg5j7qrc7v5p-8080.app.github.dev/")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// http.HandleFunc("/", handleMain)
func handleMain(w http.ResponseWriter, r *http.Request) {
    html := `<html><body><a href="/login">Login with GitHub</a></body></html>`
    fmt.Fprint(w, html)
}

func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
    url := githubOAuthConfig.AuthCodeURL(oauthStateString)
	fmt.Println("Redirect :", url)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
    if r.URL.Query().Get("state") != oauthStateString {
        http.Error(w, "Invalid state", http.StatusBadRequest)
        return
    }

    code := r.URL.Query().Get("code")
    token, err := githubOAuthConfig.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Code exchange failed", http.StatusInternalServerError)
        return
    }

    client := githubOAuthConfig.Client(context.Background(), token)
    resp, err := client.Get("https://api.github.com/user")
    if err != nil {
        http.Error(w, "Failed to get user", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    fmt.Fprintf(w, "GitHub Login Success!\n\nResponse: %s", body)
}
