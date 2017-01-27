package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
	"io/ioutil"
	"net/http"
	"os"
)

const htmlIndex = `<html><body>
<a href="/LinkedinLogin">Log in with Linkedin</a>
</body></html>
`

var (
	linkedinOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3500/LinkedinCallback",
		ClientID:     os.Getenv("L_ID"),
		ClientSecret: os.Getenv("L_SECRET"),
		Scopes:       []string{"r_basicprofile", "r_emailaddress"},
		Endpoint:     linkedin.Endpoint,
	}
	oauthStateString = "random"
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/LinkedinLogin", handleLinkedinLogin)
	http.HandleFunc("/LinkedinCallback", handleLinkedinCallback)
	fmt.Println(http.ListenAndServe(":3500", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func handleLinkedinLogin(w http.ResponseWriter, r *http.Request) {
	url := linkedinOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleLinkedinCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	token, err := linkedinOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://api.linkedin.com/v1/people/~:(id,num-connections,picture-url,location,summary,positions)?format=json&oauth2_access_token=" + token.AccessToken)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)
}
