/*
Copyright © 2024 Bridge Digital
*/
package encrypter

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/keypubfile"
	"gitea.bridge.digital/bridgedigital/db-manager-client-cli-go/services/predefined"
)

func EncryptData(dbDumpData map[string]string, keyPubFileName string) []byte {
	dbDumpDataEncoded, err := json.Marshal(dbDumpData)
	if err != nil {
		fmt.Println(predefined.BuildError("Error encoding to json:"), err)
		return nil
	}

	keyData := keypubfile.ReadKeyPubFile(keyPubFileName)
	if len(strings.TrimSpace(keyData)) == 0 {
		fmt.Println(predefined.BuildWarning("The public key is empty. Please re-create a public key."))
		return nil
	}

	pubKeyBytes := []byte(keyData)

	// Decode the PEM block
	block, _ := pem.Decode(pubKeyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		fmt.Println(predefined.BuildError("Failed to decode PEM block containing public key."))
		return nil
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Printf(predefined.BuildError("Error parsing public key: %v"), err)
		return nil
	}

	// Assert the type to *rsa.PublicKey
	rsaPubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		fmt.Println(predefined.BuildWarning("Not an RSA public key."))
		return nil
	}

	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, dbDumpDataEncoded)
	if err != nil {
		fmt.Printf(predefined.BuildError("Error encrypting data: %v"), err)
		return nil
	}

	return encryptedData
}
