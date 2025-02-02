package gategeyser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.minekube.com/gate/pkg/util/uuid"
)

const API_URL = "https://api.geysermc.org/v2/"

type GamertagXuidResult struct {
	Xuid int64 `json:"xuid"`
}

func (g *GamertagXuidResult) Uuid() (uuid.UUID, error) {
	xuid16 := strconv.FormatInt(g.Xuid, 16)

	return uuid.Parse(fmt.Sprintf("00000000-0000-0000-000%s-%s", xuid16[0:1], xuid16[1:]))
}

func GetXuid(gamertag string) (*GamertagXuidResult, error) {

	var result GamertagXuidResult
	err := geyserApiGet("xbox/xuid/"+gamertag, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

type LinkedAccountResult struct {
	BedrockID      int64     `json:"bedrock_id"`
	JavaID         uuid.UUID `json:"java_id"`
	JavaName       string    `json:"java_name"`
	LastNameUpdate int64     `json:"last_name_update"`
}

func (l *LinkedAccountResult) UnmarshalJSON(data []byte) error {
	type Alias LinkedAccountResult
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	}

	var err error

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if l.JavaID, err = uuid.Parse(aux.JavaID.String()); err != nil {
		return err
	}

	return nil
}

func GetLinkedAccount(xuid int64) (*LinkedAccountResult, error) {

	var result LinkedAccountResult

	err := geyserApiGet("link/bedrock/"+strconv.FormatInt(xuid, 10), &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

type SkinResult struct {
	Hash      string `json:"hash"`
	Steve     bool   `json:"is_steve"`
	Signature string `json:"signature"`
	TextureID string `json:"texture_id"`
	Value     string `json:"value"`
}

func GetSkin(xuid int64) (*SkinResult, error) {
	var result SkinResult
	err := geyserApiGet("skin/"+strconv.FormatInt(xuid, 10), &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func geyserApiGet(url string, result interface{}) error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", API_URL+url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("geyser api returned status code %d", res.StatusCode)
	}

	defer res.Body.Close()

	// decode to result
	if err = json.NewDecoder(res.Body).Decode(result); err != nil {
		return err
	}

	return nil
}
