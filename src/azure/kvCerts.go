package azure

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

type keyVaultCertsClient struct {
	azClient azcertificates.Client
}

var (
	certsCtx = context.Background()
)

var getKeyVaultCertsOnce sync.Once

var certsClientInstance *keyVaultCertsClient

func GetCertsClient() *keyVaultCertsClient {
	if certsClientInstance == nil {
		getKeyVaultCertsOnce.Do(
			func() {
				keyVaultName := os.Getenv("AZURE_KEY_VAULT_NAME")
				keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

				cred := GetAZIdentity()

				client, err := azcertificates.NewClient(keyVaultUrl, cred, nil)
				if err != nil {
					panic(err)
				}

				certsClientInstance = &keyVaultCertsClient{
					azClient: *client,
				}
			})
	}

	return certsClientInstance
}

func (client keyVaultCertsClient) CreateCert(certName string) {
	resp, err := client.azClient.BeginCreateCertificate(certsCtx, certName, azcertificates.CertificatePolicy{
		IssuerParameters: &azcertificates.IssuerParameters{
			Name: to.StringPtr("Self"),
		},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{
			Subject: to.StringPtr("CN=DefaultPolicy"),
		},
	}, nil)

	if err != nil {
		panic(err)
	}

	pollerResp, err := resp.PollUntilDone(certsCtx, 1*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created certificate with ID: %s\n", *pollerResp.ID)
}

func (client keyVaultCertsClient) GetCert(certName string) {
	getResp, err := client.azClient.GetCertificate(certsCtx, certName, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Cert:", *getResp.ID)
}

func (client keyVaultCertsClient) UpdateCert(certName string) {
	_, err := client.azClient.UpdateCertificateProperties(certsCtx, certName, &azcertificates.UpdateCertificatePropertiesOptions{
		Version: "newVersion",
		CertificateAttributes: &azcertificates.CertificateProperties{
			Enabled: to.BoolPtr(false),
			Expires: to.TimePtr(time.Now().Add(72 * time.Hour)),
		},
		Tags: map[string]string{"Owner": "SRE"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated certificate properites: Enabled=false, Expires=72h, Tags=SRE")
}

func (client keyVaultCertsClient) DeleteCert(certName string) {
	client.azClient.BeginDeleteCertificate(certsCtx, certName, nil)

	fmt.Println("Deleted")
}
