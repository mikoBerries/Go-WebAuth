package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	Name string
}

func main() {
	http.HandleFunc("/encode", foo)
	http.HandleFunc("/decode", bar)
	http.HandleFunc("/hash", hash)
	http.ListenAndServe(":8080", nil)

}
func hash(w http.ResponseWriter, r *http.Request) {
	//standart od school encoding using bas64
	fmt.Fprintln(w, base64.StdEncoding.EncodeToString([]byte("user:password")))
	password := "This is my password"
	bs, err := hashPassword(password)
	if err != nil {
		log.Panic(err)
	}
	err = compareHash(password, bs)
	if err != nil {
		log.Panic(err)
	}
}

func hashPassword(password string) ([]byte, error) {
	//default cost = int 10
	byteString, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error when hashpassword %w", err)
	}
	return byteString, nil
}
func compareHash(password string, bs []byte) error {
	//if returning error meanign there are not matched or something else
	err := bcrypt.CompareHashAndPassword([]byte(password), bs)
	if err != nil {
		return fmt.Errorf("Error when comparing hash %w", err)
	}
	return nil
}
func foo(w http.ResponseWriter, r *http.Request) {

	p1 := Person{
		"jamesbond",
	}
	p2 := Person{
		"bondJames",
	}
	people := []Person{p1, p2}
	err := json.NewEncoder(w).Encode(people)
	if err != nil {
		log.Println("do somethings", err)
	}
}

func bar(w http.ResponseWriter, r *http.Request) {
	p1 := Person{}
	err := json.NewEncoder(w).Encode(&p1)
	if err != nil {
		log.Println("do somethings", err)
	}
}
