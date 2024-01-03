package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

//secret key that is needed to create the JWT token is stored in a env file

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You have accessed your secret information")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(("unexpected sigining method"))
				}

				// aud := "sample.audience.io"
				// checkAudience := t.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

				// if !checkAudience {
				// 	return nil, fmt.Errorf(("invalid aud"))
				// }

				// iss := "jwtgo.io"
				// checkIss := t.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)

				// if !checkIss {
				// 	return nil, fmt.Errorf(("invalid iss"))
				// }

				// if !t.Valid {
				// 	return nil, fmt.Errorf(("invalid token"))
				// }

				return MySigningKey, nil
			})

			if err != nil {
				fmt.Fprint(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "No authorization token provided")
		}
	})
}

func handleRequest() {
	http.Handle("/", isAuthorized(homePage))
	//call the function to handle the request
	log.Fatal(http.ListenAndServe(":9001", nil))
	//starts the server in the given port no
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MySigningKey = []byte(os.Getenv("SECRET_KEY"))

	fmt.Printf("api server listening on port 9001")
	handleRequest()
}
