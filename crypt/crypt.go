package crypt

import (
	"crypto/aes"
	"io"
	"crypto/rand"
	"crypto/cipher"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha512"
	"crypto/md5"
	"errors"
)

const (
	keyLength           = 32
	keySalt             = "a simple salt for terraform"
	keyDeriveIterations = 4096
)

type Message struct {
	PlainText []byte
	/*
	Cipher Anatomy:
	[0-15]: CFB IV
	[16-31]: md5 hash
	[32->]: payload
	 */
	CipherText []byte
}

func (message *Message) Encrypt(password string) {
	block := buildBlock(password)

	// Generate a random iv and add it at the front of the ciphertext
	ciphertext := make([]byte, aes.BlockSize+md5.Size+len(message.PlainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	md5hash := md5.Sum(message.PlainText)
	for i := 0; i < md5.Size; i++ {
		ciphertext[aes.BlockSize+i] = md5hash[i]
	}

	// Encrypt the message
	cipher.NewCFBEncrypter(block, iv).
		XORKeyStream(ciphertext[aes.BlockSize+md5.Size:], message.PlainText)

	// Save changes after success
	message.CipherText = ciphertext
}

func (message *Message) Decrypt(password string) error {
	block := buildBlock(password)

	plaintext := make([]byte, len(message.CipherText)-aes.BlockSize-md5.Size)

	// Extract known md5 hash from ciphertext
	md5Hash := [16]byte{}
	for i, b := range message.CipherText[aes.BlockSize : aes.BlockSize+md5.Size] {
		md5Hash[i] = b
	}

	cipher.
		NewCFBDecrypter(block, message.CipherText[:aes.BlockSize]).
		XORKeyStream(plaintext, message.CipherText[aes.BlockSize+md5.Size:])

	// Verify results against known hash
	if md5Hash != md5.Sum(plaintext) {
		return errors.New("invalid password")
	}
	message.PlainText = plaintext
	return nil
}

func buildBlock(password string) cipher.Block {
	key := buildKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	return block
}

func buildKey(password string) []byte {
	return pbkdf2.Key([]byte(password), []byte(keySalt), keyDeriveIterations, keyLength, sha512.New)
}
