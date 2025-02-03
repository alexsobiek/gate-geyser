package floodgate

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

const IV_LENGTH = 12
const TAG_BIT_LENGTH = 128

const MAGIC = 0x3E
const SPLITTER = 0x21

const VERSION = 0
const IDENTIFIER = "^Floodgate^"

var IDENTIFIER_BYTES = []byte(IDENTIFIER)
var HEADER = IDENTIFIER + string(rune(VERSION+MAGIC))

type AesCipher struct {
	key []byte
}

func NewAesCipher(key []byte) (*AesCipher, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key length for AES")
	}

	return &AesCipher{key: key}, nil
}

func (c *AesCipher) Decrypt(cipherTextWithIv []byte) ([]byte, error) {
	if len(cipherTextWithIv) < len(HEADER)+IV_LENGTH+1 {
		return nil, errors.New("invalid ciphertext length")
	}

	headerLen := len(HEADER)
	if !bytes.HasPrefix(cipherTextWithIv, []byte(HEADER)) {
		return nil, errors.New("invalid header")
	}

	data := cipherTextWithIv[headerLen:]
	splitIndex := bytes.IndexByte(data, SPLITTER)
	if splitIndex == -1 {
		return nil, errors.New("invalid format, missing splitter")
	}

	iv := data[:splitIndex]
	cipherText := data[splitIndex+1:]

	iv, err := base64.StdEncoding.DecodeString(string(iv))
	if err != nil {
		return nil, err
	}

	cipherText, err = base64.StdEncoding.DecodeString(string(cipherText))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := gcm.Open(nil, iv, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
