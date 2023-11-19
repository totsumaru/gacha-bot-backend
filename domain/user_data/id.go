package user_data

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// IDです
//
// サーバーID+ユーザーIDを結合したものです。
type ID struct {
	value string
}

// IDを作成します
func NewID(serverID domain.DiscordID, userID domain.DiscordID) (ID, error) {
	id := ID{
		value: serverID.String() + userID.String(),
	}

	if err := id.validate(); err != nil {
		return ID{}, errors.NewError("IDの生成に失敗しました", err)
	}

	return id, nil
}

// 値を取得します
func (id ID) String() string {
	return id.value
}

// 値が空かどうかを返します
func (id ID) IsEmpty() bool {
	return id.value == ""
}

// 値を検証します
func (id ID) validate() error {
	if id.IsEmpty() {
		return errors.NewError("IDが空です", nil)
	}

	return nil
}

// 値をJSONに変換します
func (id ID) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: id.value,
	}

	return json.Marshal(data)
}

// JSONから値を復元します
func (id *ID) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("IDの復元に失敗しました", err)
	}

	id.value = data.Value

	return nil
}
