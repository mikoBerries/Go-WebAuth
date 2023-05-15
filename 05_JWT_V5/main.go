package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaim struct {
	jwt.RegisteredClaims
	UserId   string
	UserName string
	Role     string
}

var mySigningKey string

func main() {
	mySigningKey = "mySuperSecretKey"
	claims := &MyClaim{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "server",
			ID:        "1",
			Subject:   "",
			Audience:  jwt.ClaimStrings{},
		},
		"1",
		"this is username",
		"admin",
	}

	//get chiper
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenstring, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		log.Fatalln(err)
	}
	//signed string
	fmt.Println(tokenstring)
	// eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2ZXIiLCJleHAiOjE2ODQxNDk3MTUsIm5iZiI6MTY4NDE0OTExNSwiaWF0IjoxNjg0MTQ5MTE1LCJqdGkiOiIxIn0.e_42NbFu0pQ9N_UXxRAZDkurJ9zNtY9flRZCpzzlAwD9051jXzoIZiUFxh5vwQwZxV0M2cNvFTVcglEbgMnJcA
	//jwt are 3 part structure separate with "." example (Header.Payload.Signature)
	//result can check at jwt.io with "mySuperSecretKey" in verify signature
	//don't store sensitif data in jwt payload since header and payload only encoded by Base64URL

	t, err := jwt.ParseWithClaims(tokenstring, &MyClaim{},
		func(t *jwt.Token) (interface{}, error) {
			//Keyfunc will be used by the Parse methods as a callback function to supply the key for verification.
			//The function receives the parsed, but unverified Token.
			//This allows you to use properties in the Header of the token (such as `kid`) to identify which key to use.
			//vurlnable always check hash method if same that we use
			if t.Method.Alg() != jwt.SigningMethodHS512.Name {
				return nil, fmt.Errorf("eror when parse signing method")
			}
			t.Method.Verify()
			//return key to verified this token string
			return []byte(mySigningKey), nil
		})
	// somethings := t.Header["something"]

	if err != nil {
		log.Fatalln(err)
	}
	if t.Valid { //if token are valid signed using mySigningKey

		claims := t.Claims.(*MyClaim)
		fmt.Println("user ID\t", claims.UserId)
		fmt.Println("UserName\t", claims.UserName)
		fmt.Println("Role\t", claims.Role)
		// claims := t.Claims.(*jwt.MapClaims)
		// &claims["user"]
		// fmt.Println("user ID\t", claims["UserId"])
		// fmt.Println("UserName\t", claims["UserName"])
		// fmt.Println("Role\t", claims["Role"])
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		fmt.Println("That's not even a token")
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		// Invalid signature
		fmt.Println("Invalid signature")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		fmt.Println("Timing is everything")
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
}
