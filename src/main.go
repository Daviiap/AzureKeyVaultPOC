package main

import (
	"fmt"
	"kvPoc/src/kv"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func main() {
	keyVaultName := os.Getenv("AZURE_KEY_VAULT_NAME")
	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	client, err := azsecrets.NewClient(keyVaultUrl, cred, nil)
	if err != nil {
		log.Fatalf("failed to create a client: %v", err)
	}

	kv.CreateAZKeyVaultSecret(client, "quickstart-secret", "createdByGo")
	kv.CreateAZKeyVaultSecret(client, "quickstart-secret", "createdByGo2")
	kv.GetAZKeyVaultSecret(client, "quickstart-secret")
	kv.DeleteAZKeyVaultSecret(client, "quickstart-secret")
}
