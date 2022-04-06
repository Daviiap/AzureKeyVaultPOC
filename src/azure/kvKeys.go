package azure

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

type keyVaultKeysClient struct {
	azClient azkeys.Client
}

var getKeyVaultKeysOnce sync.Once

var keysClientInstance *keyVaultKeysClient

func GetKeysClient() *keyVaultKeysClient {
	if keysClientInstance == nil {
		getKeyVaultKeysOnce.Do(
			func() {
				keyVaultName := os.Getenv("AZURE_KEY_VAULT_NAME")
				keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

				cred := GetAZIdentity()

				client, err := azkeys.NewClient(keyVaultUrl, cred, nil)
				if err != nil {
					panic(err)
				}

				keysClientInstance = &keyVaultKeysClient{
					azClient: *client,
				}
			})
	}

	return keysClientInstance
}

func (client keyVaultKeysClient) CreateRSAKey(keyName string) {
	_, err := client.azClient.CreateRSAKey(context.TODO(), keyName, &azkeys.CreateRSAKeyOptions{
		Size: to.Int32Ptr(2048),
	})
	if err != nil {
		log.Fatalf("failed to create rsa key: %v", err)
	}
	fmt.Println("Created new RSA key")
}

func (client keyVaultKeysClient) UpdateKey(keyName string) {
	updateResp, err := client.azClient.UpdateKeyProperties(context.TODO(), keyName, &azkeys.UpdateKeyPropertiesOptions{
		Properties: &azkeys.Properties{
			Enabled: to.BoolPtr(false),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Key Enabled attribute set to: %t\n", *updateResp.Properties.Enabled)
}

func (client keyVaultKeysClient) DeleteKey(keyName string) {
	client.azClient.BeginDeleteKey(context.TODO(), keyName, nil)

	fmt.Println("Successfully deleted key ")
}
