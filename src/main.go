package main

import (
	"encoding/json"
	"fmt"
	"kvPoc/src/kv"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Secret *string `json:"secret"`
}

var secretName = "secret"

func HandleGetSecret(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Add("Content-Type", "json")
	client := kv.GetClient()

	secret := client.GetAZKeyVaultSecret(secretName)
	json.NewEncoder(responseWriter).Encode(Response{
		Secret: secret.Value,
	})
}

func generateRandomSecrets() {
	client := kv.GetClient()

	for {
		client.CreateAZKeyVaultSecret(secretName, fmt.Sprint(rand.Int()))
		time.Sleep(10 * time.Second)
	}
}

func main() {
	go generateRandomSecrets()

	router := mux.NewRouter()

	router.HandleFunc("/secret", HandleGetSecret).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}
