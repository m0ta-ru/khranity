package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"

	"khranity/app/utils"
)

func decrypt(cipherstring []byte, keystring string) ([]byte, error) {
	// Byte array of the string
	ciphertext := []byte(cipherstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		//panic(err)
		return nil, err
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		//panic("Text is too short")
		return nil, utils.ErrInternal
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func encrypt(plainstring []byte, keystring string) ([]byte, error) {
	// Byte array of the string
	plaintext := []byte(plainstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		//panic(err)
		return nil, err
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		//panic(err)
		return nil, err
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func writeToFile(data []byte, file string) error {
	return os.WriteFile(file, data, 0777)
}

func readFromFile(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	return data, err
}

func EncryptFile(fileIn, fileOut, token string) error {
	token = aesStringToBytes32(token)

	data, err := readFromFile(fileIn)
	if err != nil {
		return err
	}

	data, err = encrypt(data, token)
	if err != nil {
		return err
	}

	err = writeToFile(data, fileOut)
	if err != nil {
		return err
	}

	return nil
}

func DecryptFile(fileIn, fileOut, token string) error {
	token = aesStringToBytes32(token)

	data, err := readFromFile(fileIn)
	if err != nil {
		return err
	}

	data, err = decrypt(data, token)
	if err != nil {
		return err
	}

	err = writeToFile(data, fileOut)
	if err != nil {
		return err
	}

	return nil
}

func aesStringToBytes32(data string) string {
	out := [32]byte{0}
	copy(out[:], []byte(data))

	return string(out[:])
}
