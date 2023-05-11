package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

// github.com/dgrijalva/jwt-go" are depriciated
// github.com/golang-jwt/jwt/v5 succesor

//JWT  Json Web Token
// structure := {JWT standart Field}.{User Field}.Signature
type UserClaims struct {
	jwt.StandardClaims
	SessionId int64
}

func (u *UserClaims) Valid() error {
	if !u.VerifyExpiresAt(time.Now().Unix(), true) { //are token expired ?
		return fmt.Errorf("token has Expired")
	}
	if u.SessionId == 0 { //user checking
		return fmt.Errorf("invalid sessionID")
	}
	return nil
}

func createToken(c *UserClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	signedToken, err := t.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error happend in createToken when signedString %w", err)
	}
	return signedToken, nil
}

type key struct {
	key     []byte
	created time.Time
}

var keys map[string]key

func main() {

	generateNewKey()
	token, err := createToken(&UserClaims{})
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	parseToken(token)
}

func generateNewKey() error {
	newKey := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, newKey)
	if err != nil {
		return fmt.Errorf("erorr in generate new key when readfull %w", err)
	}
	uid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("erorr in generate new key when newuuidv4 %w", err)
	}

	keys[uid.String()] = key{
		key:     newKey,
		created: time.Now(),
	}

	return nil
}

func parseToken(signetdToken string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(signetdToken, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() { //checking vurnability hashing method
			return nil, fmt.Errorf("invalid signing algoritmn")
		}
		//get kid value from header
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid key id")
		}
		//check kid from database
		//exist? / expired ? / ETC
		k, ok := keys[kid]
		if !ok {
			return nil, fmt.Errorf("invalid key id")
		}
		return k.key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error in paseToken,Token is no valid")
	}
	return t.Claims.(*UserClaims), nil
}
