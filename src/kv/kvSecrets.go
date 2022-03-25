package kv

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func CreateAZKeyVaultSecret(client *azsecrets.Client, secretName string, secretValue string) azsecrets.SetSecretResponse {
	resp, err := client.SetSecret(context.TODO(), secretName, secretValue, nil)
	if err != nil {
		log.Fatalf("failed to create a secret: %v", err)
	}

	fmt.Printf("Name: %s, Value: %s\n", *resp.ID, *resp.Value)

	return resp
}

func GetAZKeyVaultSecret(client *azsecrets.Client, secretName string) azsecrets.GetSecretResponse {
	getResp, err := client.GetSecret(context.TODO(), "quickstart-secret", nil)
	if err != nil {
		log.Fatalf("failed to get the secret: %v", err)
	}

	fmt.Printf("secretValue: %s\n", *getResp.Value)

	return getResp
}

func DeleteAZKeyVaultSecret(client *azsecrets.Client, secretName string) azsecrets.DeleteSecretResponse {
	respDel, err := client.BeginDeleteSecret(context.TODO(), secretName, nil)

	if err != nil {
		log.Fatalf("failed to delete secret: %v", err)
	}

	response, err := respDel.PollUntilDone(context.TODO(), time.Second)

	if err != nil {
		log.Fatalf("failed to delete secret: %v", err)
	}

	fmt.Println(secretName + " has been deleted\n")

	return response
}
