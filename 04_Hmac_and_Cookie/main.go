package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {

	http.HandleFunc("/setcookie", setCookie)
	http.HandleFunc("/checkcookie", checkCookie)
	http.ListenAndServe(":8080", nil)
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("cookie") == "" { //no cookies
		cokievalue, err := writeCode("yourEmail@yahoo.com")
		if err != nil {
			log.Fatalln(err)
			http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		}
		c := http.Cookie{
			Name:  "session",
			Value: cokievalue + "|" + "yourEmail@yahoo.com",
		}
		http.SetCookie(w, &c)
	}

}
func checkCookie(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/setCookie", http.StatusNotFound)
	}
	ss := strings.SplitN(c.Value, "|", 0)
	if len(ss) != 2 {
		log.Fatal("error when split")
		http.Redirect(w, r, "/setCookie", http.StatusNotFound)
	}
	code := ss[0]
	email := ss[1]
	hashEmail, err := writeCode(email)
	if err != nil {
		http.Redirect(w, r, "/setCookie", http.StatusNotFound)
	}
	//comapring email and hash code we send to client as cookie
	equal := hmac.Equal([]byte(code), []byte(hashEmail))
	if !equal {
		http.Redirect(w, r, "/setCookie", http.StatusForbidden)
	}
	//pass
}
func writeCode(cookie string) (string, error) {
	//Hmac using sha 512 hash algorithm an "secret key" as key
	h := hmac.New(sha512.New, []byte("secret key"))
	_, err := h.Write([]byte(cookie))
	if err != nil {
		return "", fmt.Errorf("error hen write hmac %w", err)
	}
	result := h.Sum(nil)
	return string(result), nil
}
