# Key Vault Go Example

A simple application using [Azure Key Vault](https://azure.microsoft.com/pt-br/services/key-vault/) secrets.

## How it works

This application runs a [go routine](https://go.dev/tour/concurrency/1) to generate random secrets and saves it in a Key Vault.
And gives a simple api with one route: `GET` => `/secret`, that returns the secret value stored in Key Vault.

## QuickStart

### 1. Creating KeyVault

First, you must be logged in with Azure CLI:

```sh
az login
```

Then, you must create a Resource Group and a KeyVault instance:

```sh
az group create --name <resource_group_name> --location eastus
az keyvault create --name <key_vault_name> --resource-group <resource_group_name>
```

### 2. Setting Environment Variables

After creating the KeyVault in Azure, you must set the environment variables to run the application:

| Variable             | Description                                 |
| -------------------- | ------------------------------------------- |
| AZURE_KEY_VAULT_NAME | The name of the KeyVault that you'v created |
|   KEY_REFRESH_RATE   |     The refresh rate for the saved token    |

### 3. Running The Application

After setting everything in steps `1` and `2`, you just have to run the following commands:

```sh
go mod download
go run src/main.go
```
