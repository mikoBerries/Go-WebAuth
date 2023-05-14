package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	msg := "Hi There this is massage to send to your email enjoy your days"
	encodeBase64 := encodeBase64(msg)
	fmt.Println(encodeBase64)

	decodeBase64, err := decodeBase64(encodeBase64)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(decodeBase64)

	symetricalKey := "this is my key to encrypt and decript"
	//key are [16]byte
	bs, err := bcrypt.GenerateFromPassword([]byte(symetricalKey), bcrypt.MinCost)
	if err != nil {
		log.Fatalln("error when bycrpt password", err)
	}
	result, err := enDecode(msg, bs[:16])
	if err != nil {
		log.Fatalln(err)
	}
	//encript with CTR chiper
	fmt.Println(string(result))

	result2, err := enDecode(string(result), bs[:16])
	if err != nil {
		log.Fatalln(err)
	}
	//decript with CTR chiper
	fmt.Println(string(result2))
	//sha256 or 512 better
	fileName := "sample.txt"
	hashed, err := sha(fileName)
	if err != nil {
		// log.Panic(err) will give stack of func to terminal
		log.Fatalln(err)
	}
	fmt.Println(string(hashed))
}

func encodeBase64(msg string) string {
	return base64.URLEncoding.EncodeToString([]byte(msg))

}

func decodeBase64(encoded string) (string, error) {
	result, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("Error happend wehn decode base 64 %w", err)
	}
	return string(result), nil
}

func enDecode(msg string, key []byte) ([]byte, error) {
	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error happend when encode new chiper %w", err)
	}
	iv := make([]byte, aes.BlockSize)
	//cipher need chiperblock with key and blok size
	//NewCTR are cipher that useCTR algoritmh
	cs := cipher.NewCTR(cb, iv)
	buff := &bytes.Buffer{}
	//steam writer need chiper and addres buffer(addres to write it)
	sw := cipher.StreamWriter{
		S: cs,
		W: buff,
	}
	_, err = sw.Write([]byte(msg))
	//result in sreamWriter.W
	if err != nil {
		return nil, fmt.Errorf("error happed when stream writing %w", err)
	}
	return buff.Bytes(), nil
}

func sha(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error when opening file %w", err)
	}
	defer f.Close()

	//hash.hash implement reader
	h := sha512.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return nil, fmt.Errorf("error while copy %w", err)
	}
	resutl := h.Sum(nil)
	return resutl, nil
}
