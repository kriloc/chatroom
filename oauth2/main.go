package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

//const clientID = "<your client id>"
const clientID = "00c35de57c6ce6843e50"

//const clientSecret = "<your client secret>"
const clientSecret = "b23e09c3304c3170b8235933cb4f8cffa8715cfb"

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("oauth2/hello.html")
		t.Execute(w, nil)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("oauth2/login.html")
		t.Execute(w, nil)
	}
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/", hello)
	http.HandleFunc("/hello", hello)

	httpClient := http.Client{}
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?"+
			"client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		req.Header.Set("accept", "application/json")

		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer res.Body.Close()

		var t AccessTokenResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Location", "/hello.html?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	http.ListenAndServe(":8087", nil)
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
