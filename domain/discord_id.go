package domain

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// DiscordのIDです
type DiscordID struct {
	value string
}

// DiscordのIDを生成します
func NewDiscordID(value string) (DiscordID, error) {
	d := DiscordID{value: value}

	if err := d.validate(); err != nil {
		return d, errors.NewError("検証に失敗しました", err)
	}

	return d, nil
}

// DiscordのIDを返します
func (d DiscordID) String() string {
	return d.value
}

// DiscordのIDが存在しているか確認します
func (d DiscordID) IsEmpty() bool {
	return d.value == ""
}

// DiscordのIDを検証します
func (d DiscordID) validate() error {
	if d.value == "" {
		return errors.NewError("DiscordIDが空です")
	}

	return nil
}

// DiscordのIDをJSONに変換します
func (d DiscordID) MarshalJSON() ([]byte, error) {
	data := struct {
		DiscordID string `json:"discord_id"`
	}{
		DiscordID: d.value,
	}

	return json.Marshal(data)
}

// JSONからDiscordのIDを復元します
func (d *DiscordID) UnmarshalJSON(b []byte) error {
	data := struct {
		DiscordID string `json:"discord_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからDiscordIDの復元に失敗しました", err)
	}

	d.value = data.DiscordID

	return nil
}
