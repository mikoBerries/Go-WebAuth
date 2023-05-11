package main

import "github.com/golang-jwt/jwt/v5"

type MyClaim struct {
	jwt.MapClaims
	sessionId string
}

func main() {

}
