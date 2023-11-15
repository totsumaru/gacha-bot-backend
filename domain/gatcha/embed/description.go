package embed

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// 内容です
type Description struct {
	value string
}

// 内容を生成します
func NewDescription(value string) (Description, error) {
	d := Description{value: value}

	if err := d.validate(); err != nil {
		return Description{}, errors.NewError("内容が不正です", err)
	}

	return d, nil
}

// 内容を返します
func (d Description) String() string {
	return d.value
}

// 内容が存在しているか確認します
func (d Description) IsEmpty() bool {
	return d.value == ""
}

// 内容を検証します
func (d Description) validate() error {
	if len([]rune(d.value)) > 1500 {
		return errors.NewError("送信内容の最大文字数を超えています")
	}

	return nil
}

// 内容をJSONに変換します
func (d Description) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: d.value,
	}

	return json.Marshal(data)
}

// 内容をJSONから変換します
func (d *Description) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	d.value = data.Value

	return nil
}
