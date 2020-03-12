package sm2enc

import (
	"fmt"

	"github.com/tjfoc/gmsm/sm2"
)

//GenSM2Key generates SM2 private and public keys.
func GenSM2Key() (prvkey, pubkey []byte) {
	private, e := sm2.GenerateKey()
	if e != nil {
		fmt.Println("sm2 key generation failed！")
	}
	public := &private.PublicKey
	sm2.WritePrivateKeytoPem("private.pem", private, nil)
	//sm2.WritePublicKeytoPem("public.pem", public, nil)
	prvkey, _ = sm2.WritePrivateKeytoMem(private, nil)
	pubkey, _ = sm2.WritePublicKeytoMem(public, nil)
	return
}

//SM2Encrypt encrypts data using public key.
func SM2Encrypt(data, keyBytes []byte) []byte {
	public, _ := sm2.ReadPublicKeyFromMem(keyBytes, nil)
	bytes, i := public.Encrypt(data)
	if i != nil {
		fmt.Println("encryption failed！")
	}
	return bytes
}

//SM2Decrypt decrypts ciphertext with private key.
func SM2Decrypt(ciphertext, keyBytes []byte) []byte {
	private, _ := sm2.ReadPrivateKeyFromMem(keyBytes, nil)
	decrypt, i := private.Decrypt(ciphertext)
	if i != nil {
		fmt.Println("decryption failed！")
	}
	return decrypt
}
