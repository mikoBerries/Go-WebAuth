package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
)

var key []byte

func main() {
	for i := 1; i <= 64; i++ {
		key = append(key, byte(i))
	}
	checkSignature()
}

func signMessage(mssg []byte) ([]byte, error) {
	h := hmac.New(sha512.New, key)
	_, err := h.Write(mssg)
	if err != nil {
		return nil, fmt.Errorf("Error in Sign Message while hassing massage : %w", err)
	}
	signature := h.Sum(nil)
	return signature, nil
}
func checkSignature(msg, signatureToCheck []byte) (bool, error) {
	signature, err := signMessage(msg)
	if err != nil {
		return false, fmt.Errorf("Error in checkSignature while getting signature massage : %w", err)
	}
	result := hmac.Equal(signature, signatureToCheck)
	return result, nil
}
