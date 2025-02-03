package floodgate

import (
	"fmt"
	"strings"
)

type Floodgate struct {
	cipher *AesCipher
}

func NewFloodgate(key []byte) (*Floodgate, error) {
	cipher, err := NewAesCipher(key)
	if err != nil {
		return nil, err
	}

	return &Floodgate{cipher: cipher}, nil
}

func (f *Floodgate) Decrypt(data []byte) ([]byte, error) {
	return f.cipher.Decrypt(data)
}

func (f *Floodgate) ReadHostname(hostname string) (string, *BedrockData, error) {
	parts := strings.Split(hostname, "\u0000")

	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid hostname format: %s", hostname)
	}

	originalHostname := parts[0]

	// check if port is appended
	data := parts[1]

	if strings.Contains(data, ":") {
		data = strings.Split(data, ":")[0]
	}

	bedrockDataBytes, err := f.Decrypt([]byte(data))

	if err != nil {
		return "", nil, err
	}

	bedrockData, err := ReadBedrockData(string(bedrockDataBytes))

	if err != nil {
		return "", nil, err
	}

	return originalHostname, bedrockData, nil
}
