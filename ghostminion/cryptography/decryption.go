package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"ghostminion/config"
)

func DecryptData(cipherText []byte) ([]byte, error) {
	configInstance, err := config.GetConfig()
	if err != nil {
		return []byte{}, err
	}
	key, err := hex.DecodeString(configInstance.Installation.AESKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("cipherText too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
