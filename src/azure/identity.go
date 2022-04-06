package azure

import (
	"log"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var generateIdentityOnce sync.Once

var identity *azidentity.DefaultAzureCredential

func GetAZIdentity() *azidentity.DefaultAzureCredential {
	if identity == nil {
		generateIdentityOnce.Do(
			func() {
				var err error
				identity, err = azidentity.NewDefaultAzureCredential(nil)

				if err != nil {
					log.Fatalf("Could not create a az identity %v", err)
				}
			})
	}

	return identity
}
