package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var oc oauth2.Config

type ghUser struct {
	Data struct {
		Viewer struct {
			Id   string `json:"id"`
			Bio  string `json:"bio"`
			Name string `json:"name"`
		} `json:"viewer"`
	} `json:"data"`
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/oauth2/github", oauth2Github)
	http.HandleFunc("/recieveOauth2", recieve)

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<form action="/oauth2/github" method="POST">
			<input type="submit" value="Login with Github">
		</form>
	</body>
	</html>`)
}

func oauth2Github(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	oc = oauth2.Config{
		ClientID:     "9932db92d15ec974fead",                     //received when you first created the application
		ClientSecret: "294f1890438abfdaea0d97c2fd908dd262e133ed", //from oauth2 client server
		Endpoint:     github.Endpoint,
		// oauth2.Endpoint{ //end point to oaauth2 use to verify /exchange data
		// 	AuthURL:  "",
		// 	TokenURL: "",
		// },
		// RedirectURL: "http://localhost:8080/recieveOauth2", //an URL to redirect user after oauth2 are done some client not requer redirecturl
		Scopes: []string{"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"}, //things that we need / request to client server exp: photo,data,etc
		// Scopes: []string{"user"},
	}
	fmt.Println(oc)
	randomState := "1203123hasd" //serve as state as id so we know whos data we throw and recieve

	redirectURL := oc.AuthCodeURL(randomState)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
func recieve(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, r)
	//we recieve code and state(from us)
	//http://localhost:8080/recieveOauth2?code=77fd98f18f7a3de86693&state=1203123hasd
	code := r.FormValue("code")
	state := r.FormValue("state")

	if state != "1203123hasd" || code == "" { //store from database
		http.Error(w, "incorect state/code", http.StatusBadRequest)
		return
	}
	//exchange code to token
	token, err := oc.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Couldn't login", http.StatusInternalServerError)
		return
	}
	// ts := oc.TokenSource(r.Context(), token)
	// client := oauth2.NewClient(r.Context(), ts)
	client := oc.Client(r.Context(), token)
	//https://docs.github.com/en/graphql/guides/forming-calls-with-graphql
	//documentation for github
	rBody := strings.NewReader(`{"query":"query{viewer{id bio name}}"}`) //graphql query
	resp, err := client.Post("https://api.github.com/graphql", "application/json", rBody)
	if err != nil {
		http.Error(w, "coulnot get user data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	//io util real all respond.body to byte

	//using read all
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "couldn't read github respons body", http.StatusInternalServerError)
		return
	}
	log.Println(string(bs))
	//{"data":{"viewer":{"id":"MDQ6VXNlcjU2NDA2NjE3"}}}
	ug := ghUser{}
	json.Unmarshal(bs, &ug)

	log.Printf("unmarshaled data %#v", ug)

	//using json decoder consme io reader
	json.NewDecoder(resp.Body).Decode(&ug)

	//with data from ug we can create ne user to database
	//and following login using jwt with some uniq data from user db
}
