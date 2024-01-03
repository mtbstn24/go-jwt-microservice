package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

func GetJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "TokenCreator-101"
	claims["aud"] = "sample.audience.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Local().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(MySigningKey)

	if err != nil {
		fmt.Printf("something went wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	fmt.Println(validToken)
	if err != nil {
		fmt.Println("Failed to get the Token")
	}
	fmt.Fprint(w, string(validToken))
}

func handleRequests() {
	http.HandleFunc("/", Index)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MySigningKey = []byte(os.Getenv("SECRET_KEY"))

	handleRequests()
}
