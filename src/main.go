package main

import (
	"encoding/json"
	"fmt"
	"kvPoc/src/azure"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type SecretsResponse struct {
	Secret *string `json:"secret"`
}

var secretName = "secret"

func HandleGetSecret(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Add("Content-Type", "json")
	client := azure.GetSecretsClient()

	secret := client.GetAZKeyVaultSecret(secretName)
	json.NewEncoder(responseWriter).Encode(SecretsResponse{
		Secret: secret.Value,
	})
}

func generateRandomSecretsRoutine() {
	client := azure.GetSecretsClient()

	for {
		client.CreateAZKeyVaultSecret(secretName, fmt.Sprint(rand.Int()))

		sleepTime, parseError := strconv.Atoi(os.Getenv("KEY_REFRESH_RATE"))

		if parseError != nil {
			log.Fatal("INVALID REFRASH RATE TIME")
		}

		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func main() {
	go generateRandomSecretsRoutine()

	router := mux.NewRouter()

	router.HandleFunc("/secret", HandleGetSecret).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))

	// certsclient := azure.GetCertsClient()
	// certsclient.CreateCert("newCert")
	// certsclient.GetCert("newCert")
	// certsclient.UpdateCert("newCert")
	// certsclient.DeleteCert("newCert")
}
