package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"

	"github.com/joho/godotenv"
)


var secretKey string

func GenerateHMAC(data string)string{
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal("Error loading .env file")
	}
	secretKey=os.Getenv("SECRET_KEY")
	h:=hmac.New(sha256.New,[]byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMAC(data,signature string)bool{
	expected:=GenerateHMAC(data)
	return hmac.Equal([]byte(expected),[]byte(signature))
}