package floodgate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.minekube.com/gate/pkg/util/uuid"
)

type BedrockData struct {
	Version      string
	Username     string
	Xuid         int64
	DeviceOS     int
	Language     string
	UIProfile    int
	InputMode    int
	IP           string
	LinkedPlayer string
	Proxy        bool
	SubscribeID  string
	VerifyCode   string
}

func ReadBedrockData(data string) (*BedrockData, error) {
	parts := strings.Split(data, "\u0000")

	if len(parts) != 12 {
		return nil, errors.New("invalid bedrock data")
	}

	version := parts[0]
	username := parts[1]

	xuid, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, err
	}

	deviceOS, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, err
	}

	language := parts[4]

	uiProfile, err := strconv.Atoi(parts[5])
	if err != nil {
		return nil, err
	}

	inputMode, err := strconv.Atoi(parts[6])
	if err != nil {
		return nil, err
	}

	ip := parts[7]
	linkedPlayer := parts[8]
	proxy := parts[9] == "1"
	subscribeID := parts[10]
	verifyCode := parts[11]

	return &BedrockData{
		Version:      version,
		Username:     username,
		Xuid:         xuid,
		DeviceOS:     deviceOS,
		Language:     language,
		UIProfile:    uiProfile,
		InputMode:    inputMode,
		IP:           ip,
		LinkedPlayer: linkedPlayer,
		Proxy:        proxy,
		SubscribeID:  subscribeID,
		VerifyCode:   verifyCode,
	}, nil
}

func (d *BedrockData) JavaUuid() (uuid.UUID, error) {
	xuid16 := strconv.FormatInt(d.Xuid, 16)

	return uuid.Parse(fmt.Sprintf("00000000-0000-0000-000%s-%s", xuid16[0:1], xuid16[1:]))
}
