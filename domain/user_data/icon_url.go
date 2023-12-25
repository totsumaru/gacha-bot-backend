package user_data

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// アイコンURLです
type IconURL struct {
	value string
}

// アイコンURLを作成します
func NewIconURL(value string) (IconURL, error) {
	i := IconURL{
		value: value,
	}

	if err := i.validate(); err != nil {
		return IconURL{}, errors.NewError("アイコンURLの生成に失敗しました", err)
	}

	return i, nil
}

// 値を取得します
func (i IconURL) String() string {
	return i.value
}

// 値が空かどうかを返します
func (i IconURL) IsEmpty() bool {
	return i.value == ""
}

// 値を検証します
func (i IconURL) validate() error {
	return nil
}

// 値をJSONに変換します
func (i IconURL) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: i.value,
	}

	return json.Marshal(data)
}

// JSONから値を復元します
func (i *IconURL) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("IconURLの復元に失敗しました", err)
	}

	i.value = data.Value

	return nil
}
