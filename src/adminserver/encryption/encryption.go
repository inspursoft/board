package encryption

import (
	"github.com/inspursoft/board/src/adminserver/rsaenc"
	"github.com/inspursoft/board/src/adminserver/sm2enc"
	"os"
)

//GenKey generates public and private keys.
func GenKey(algo string) (prvkey, pubkey []byte) {
	if algo == "rsa" {
		prvkey, pubkey = rsaenc.GenRsaKey()
	} else {
		prvkey, pubkey = sm2enc.GenSM2Key()
	}
	return
}

//Encrypt using selected algorithm.
func Encrypt(algo string, data []byte, keyBytes []byte) []byte {
	var ciphertext []byte
	if algo == "rsa" {
		ciphertext = rsaenc.RsaEncrypt(data, keyBytes)
	} else {
		ciphertext = sm2enc.SM2Encrypt(data, keyBytes)
	}
	return ciphertext
}

//Decrypt using specified algorithm.
func Decrypt(algo string, ciphertext []byte, keyBytes []byte) []byte {
	var password []byte
	if algo == "rsa" {
		password = rsaenc.RsaDecrypt(ciphertext, keyBytes)
	} else {
		password = sm2enc.SM2Decrypt(ciphertext, keyBytes)
	}
	return password
}

//CheckFileIsExist checks the fire exist or not.
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
