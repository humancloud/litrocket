package dataencry

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"litrocket/common"
)

// Encrypt And Marshal JSON.
func Marshal(v interface{}) ([]byte, error) {
	buf, err := json.Marshal(v)
	return buf, err
}

// Decrypt And Unmarshal JSON.
func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	return err
}

// Encrypt Users Password Save To Database.
func EncryptUserPasswd() {

}

// Decrypt Users Password From Database.
func DecryptUserPasswd() {

}

// Encrypt Aes Passwd From Database.
func EncryptPasswd(pass string) string {
	origData := []byte(pass)     // 待加密的数据
	key := []byte(common.AESKEY) // 加密的密钥
	encrypted := AesEncryptECB(origData, key)
	return base64.StdEncoding.EncodeToString(encrypted)
}

// Decrypt Aes Passwd From Network.  解密使用AES ECB PKCS7Padding.
func DecryptPasswd(pass string) string {
	m, _ := base64.StdEncoding.DecodeString(pass)
	encrypted := []byte(m)
	key := []byte(common.AESKEY) // 加密的密钥
	decrypted := AesDecryptECB(encrypted, key)
	return string(decrypted)
}

// Aes,Ecb,pkcs7padding加密
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
