package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	clientPrivateKey *rsa.PrivateKey
	clientPublicKey  *rsa.PublicKey
)

const (
	clientPrivateFile = "client_pvt.pem"
	clientPubFile     = "client_pub.pem"
	keySize           = 2048
)

func init() {
	if err := generateAndLoadKeys(); err != nil {
		log.Fatalf("Failed to generate/load keys: %v", err)
	}
}

func generateAndLoadKeys() error {
	if !fileExists(clientPrivateFile) || !fileExists(clientPubFile) {
		fmt.Println("Key files do not exist. Generating new keys...")
		if err := generatePEMFiles(); err != nil {
			return fmt.Errorf("failed to generate PEM files: %v", err)
		}
	} else {

		fmt.Println("Both key files exist.Skipping key generation...")

	}

	var err error
	clientPrivateKey, clientPublicKey, err = loadClientKeys(clientPrivateFile, clientPubFile)
	if err != nil {
		return fmt.Errorf("failed to load client keys: %v", err)
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func generatePEMFiles() error {
	key, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}

	if err := savePEMKey(clientPrivateFile, "PRIVATE KEY", x509.MarshalPKCS8PrivateKey, key); err != nil {
		return err
	}
	if err := savePEMKey(clientPubFile, "PUBLIC KEY", x509.MarshalPKIXPublicKey, &key.PublicKey); err != nil {
		return err
	}

	return nil
}

func savePEMKey(filename, keyType string, marshalFunc func(interface{}) ([]byte, error), key interface{}) error {
	keyBytes, err := marshalFunc(key)
	if err != nil {
		return err
	}
	keyBlock := &pem.Block{
		Type:  keyType,
		Bytes: keyBytes,
	}
	keyFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	if err := pem.Encode(keyFile, keyBlock); err != nil {
		return err
	}
	return nil
}

func loadClientKeys(privateKeyPath, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := loadPublicKey(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, publicKey, nil
}

func loadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyData)
	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	return privateKey.(*rsa.PrivateKey), nil
}

func loadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	publicKeyData, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %v", err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyData)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	return publicKeyInterface.(*rsa.PublicKey), nil
}

func GetClientPrivateKey() *rsa.PrivateKey {
	return clientPrivateKey
}

func GetClientPublicKey() *rsa.PublicKey {
	return clientPublicKey
}
