package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
)

func key(k string) string {
	for len(k) < 32 {
		k += base32.StdEncoding.EncodeToString([]byte(k))
	}
	return k[:32]
}

func Encrypt(secret string, plaintext string) (string, error) {
	blocks, err := aes.NewCipher([]byte(key(secret)))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blocks)
	if err != nil {
		return "", err
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

func Decrypt(secret string, ciphertext string) (string, error) {
	blocks, err := aes.NewCipher([]byte(key(secret)))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blocks)
	if err != nil {
		return "", err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
