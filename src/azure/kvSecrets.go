package azure

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

type keyVaultSecretsClient struct {
	azClient azsecrets.Client
}

var getSecretsClientOnce sync.Once

var secretsClientInstance *keyVaultSecretsClient

func GetSecretsClient() *keyVaultSecretsClient {
	if secretsClientInstance == nil {
		getSecretsClientOnce.Do(
			func() {
				keyVaultName := os.Getenv("AZURE_KEY_VAULT_NAME")
				keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

				cred := GetAZIdentity()

				client, err := azsecrets.NewClient(keyVaultUrl, cred, nil)
				if err != nil {
					log.Fatalf("failed to create a client: %v", err)
				}

				secretsClientInstance = &keyVaultSecretsClient{
					azClient: *client,
				}
			})
	}

	return secretsClientInstance
}

func (client keyVaultSecretsClient) CreateAZKeyVaultSecret(secretName string, secretValue string) azsecrets.SetSecretResponse {
	resp, err := client.azClient.SetSecret(context.TODO(), secretName, secretValue, nil)
	if err != nil {
		log.Fatalf("failed to create a secret: %v", err)
	}

	fmt.Printf("Name: %s, Value: %s\n", *resp.ID, *resp.Value)

	return resp
}

func (client keyVaultSecretsClient) GetAZKeyVaultSecret(secretName string) azsecrets.GetSecretResponse {
	getResp, err := client.azClient.GetSecret(context.TODO(), secretName, nil)
	if err != nil {
		log.Fatalf("failed to get the secret: %v", err)
	}

	fmt.Printf("secretValue: %s\n", *getResp.Value)

	return getResp
}

func (client keyVaultSecretsClient) DeleteAZKeyVaultSecret(secretName string) bool {
	respDel, _ := client.azClient.BeginDeleteSecret(context.TODO(), secretName, nil)

	for respDel.Poller.Done() {
	}

	fmt.Println(secretName + " has been deleted\n")

	return respDel.Poller.Done()
}
