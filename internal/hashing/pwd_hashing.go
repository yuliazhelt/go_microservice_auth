package hashing

import (
	"bytes"
	"encoding/base64"
	
	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"
)

var iter = 4096
var keyLen = 32

func HashPassword(password string, saltString string) string {
	salt := bytes.NewBufferString(saltString).Bytes()
	df := pbkdf2.Key([]byte(password), salt, iter, keyLen, sha1.New)
	cipherText := base64.StdEncoding.EncodeToString(df)
	return cipherText
}

func VerifyPassword(password string, cipherText string, saltString string) bool {
	salt := bytes.NewBufferString(saltString).Bytes()
	df := pbkdf2.Key([]byte(password), salt, iter, keyLen, sha1.New)
	return equal(cipherText, df)
}

func equal(cipherText string, newCipherText []byte) bool {
	x, _ := base64.StdEncoding.DecodeString(cipherText)
	diff := uint64(len(x)) ^ uint64(len(newCipherText))
	for i := 0; i < len(x) && i < len(newCipherText); i++ {
		diff |= uint64(x[i]) ^ uint64(newCipherText[i])
	}
	return diff == 0
}
