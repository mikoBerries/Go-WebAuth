package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//fake DB schema
type user struct {
	password []byte
	First    string
}

type MyClaim struct {
	jwt.RegisteredClaims
	SID string
}

var db = map[string]user{}
var secretKey = []byte("this is secret")
var sessions = map[string]string{}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)

	http.HandleFunc("/login", login)

	http.ListenAndServe(":8080", nil)
}
func index(w http.ResponseWriter, r *http.Request) {
	//check cookie sessionId
	c, err := r.Cookie("sessionID")
	if err != nil {
		c = &http.Cookie{
			Name:  "sessionID",
			Value: "",
		}
	}

	s, err := parseToken(c.Value)
	if err != nil {
		log.Println("index parseToken", err)
	}

	var e string
	if s != "" {
		e = sessions[s]
	}
	var f string
	if user, ok := db[e]; ok {
		f = user.First
	}

	errMsg := r.FormValue("msg")

	fmt.Fprintf(w, `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Document</title>
	</head>
	<body>
	<h1>IF YOU HAVE A SESSION, HERE IS YOUR NAME: %s</h1>
	<h1>IF YOU HAVE A SESSION, HERE IS YOUR EMAIL: %s</h1>
	<h1>IF THERE IS ANY MESSAGE FOR YOU, HERE IT IS: %s</h1>
        <h1>REGISTER</h1>
		<form action="/register" method="POST">
		<label for="first">First</label>
		<input type="text" name="first" placeholder="First" id="first">
		<input type="email" name="e">
			<input type="password" name="p">
			<input type="submit">
        </form>
        <h1>LOG IN</h1>
        <form action="/login" method="POST">
            <input type="email" name="e">
			<input type="password" name="p">
			<input type="submit">
        </form>
	</body>
	</html>`, f, e, errMsg)
}
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	e := r.FormValue("e")
	if e == "" {
		msg := url.QueryEscape("your email needs to not be empty")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	p := r.FormValue("p")
	if p == "" {
		msg := url.QueryEscape("your email password needs to not be empty")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	f := r.FormValue("first")
	if p == "" {
		msg := url.QueryEscape("first name needs to not be empty")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	bsp, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		msg := "there was an internal server error - evil laugh: hahahahaha"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	log.Println("password", e)
	log.Println("bcrypted", bsp)

	db[e] = user{
		password: bsp,
		First:    f,
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := url.QueryEscape("your method was not post")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	e := r.FormValue("e")
	if e == "" {
		msg := url.QueryEscape("your email needs to not be empty")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	p := r.FormValue("p")
	if p == "" {
		msg := url.QueryEscape("your email password needs to not be empty")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	if _, ok := db[e]; !ok {
		msg := url.QueryEscape("your email or password didn't match")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	//compare Hash in db and incoming
	err := bcrypt.CompareHashAndPassword(db[e].password, []byte(p))
	if err != nil {
		msg := url.QueryEscape("your email or password didn't match")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	sUUID, err := uuid.NewV7()
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	sessions[sUUID.String()] = e
	token, err := createToken(sUUID.String())
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	c := http.Cookie{
		Name:  "sessionID",
		Value: token,
	}
	http.SetCookie(w, &c)
	msg := url.QueryEscape("you logged in " + e)
	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
}

func createToken(sessionId string) (string, error) {
	//Hmac encode
	/*
		h := hmac.New(sha512.New, secretKey)
		_, err := h.Write([]byte(sessionId))
		if err != nil {
			return "", fmt.Errorf("eror when create token hash write %w", err)
		}
		result := base64.StdEncoding.EncodeToString(h.Sum(nil))

		return result + "|" + sessionId, nil
	*/

	//jwt encode
	claim := &MyClaim{
		jwt.RegisteredClaims{
			Issuer:    "localhost:8080",
			Audience:  jwt.ClaimStrings{},
			Subject:   "session",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "1",
		},
		sessionId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("eror when createToken signed %w", err)
	}

	return ss, nil
}

func parseToken(token string) (string, error) {
	//Hmac
	/*
		xs := strings.SplitN(token, "|", 2)
		if len(xs) != 2 {
			return "", fmt.Errorf("null string")
		}
		signedSessionId := xs[0]
		bs, err := base64.StdEncoding.DecodeString(signedSessionId)
		if err != nil {
			return "", fmt.Errorf("error when parseToken - decode %w", err)
		}

		h := hmac.New(sha512.New, secretKey)
		_, err = h.Write([]byte(xs[1]))
		if err != nil {
			return "", fmt.Errorf("error when parseToken - write %w", err)
		}
		ok := hmac.Equal([]byte(bs), h.Sum(nil))
		if !ok {
			return "", fmt.Errorf("couldn't parseToken not equal signed sid and sid")
		}

		return xs[1], nil
	*/

	//JWT
	t, err := jwt.ParseWithClaims(token, &MyClaim{}, func(t *jwt.Token) (interface{}, error) {
		//alwas check header hash method
		if t.Method.Alg() != jwt.SigningMethodHS512.Name {
			return "", fmt.Errorf("error parseToken SigningMethod not allowed")
		}
		return secretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("error parseToken parse with claim %w", err)
	}

	if !t.Valid { // not valid signature with provided secret key
		if errors.Is(err, jwt.ErrTokenMalformed) {
			err = fmt.Errorf("that's not even a token")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			// Invalid signature
			err = fmt.Errorf("invalid signature")
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			// Token is either expired or not active yet
			err = fmt.Errorf("token is expired")
		} else {
			err = fmt.Errorf("couldn't handle this token: %w", err)
		}

		return "", err
	}
	c := t.Claims.(*MyClaim)
	if c.SID == "" {
		return "", errors.New("empty sid")
	}

	return c.SID, nil

}
