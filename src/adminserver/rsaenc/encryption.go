package rsaenc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

//GenRsaKey generates RSA public and private keys.
func GenRsaKey() (prvkey, pubkey []byte) {
	// generate private key file
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	file, err := os.Create("private.pem")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	err = pem.Encode(file, block)

	//return []byte type private and public keys.
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}

//RsaSignWithSha256 signs data with private key.
func RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

//RsaVerySignWithSha256 validates data using public key.
func RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

//RsaEncrypt using public key.
func RsaEncrypt(data, keyBytes []byte) []byte {
	//decode pem public key
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	//parse the key
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//type assertion
	pub := pubInterface.(*rsa.PublicKey)
	//encrypting
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

//RsaDecrypt using private key.
func RsaDecrypt(ciphertext, keyBytes []byte) []byte {
	//acquire private key
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	//parse PKCS1 private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//decrypting
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		panic(err)
	}
	return data
}
