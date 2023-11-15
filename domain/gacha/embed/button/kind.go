package button

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

const (
	ButtonKindToOpen   = "to_open"
	ButtonKindToResult = "to_result"
)

// ボタンの種類です
type Kind struct {
	value string
}

// ボタンの種類を生成します
func NewKind(value string) (Kind, error) {
	t := Kind{value: value}

	if err := t.validate(); err != nil {
		return Kind{}, errors.NewError("ボタンの種類が不正です", err)
	}

	return t, nil
}

// ボタンの種類を返します
func (k Kind) String() string {
	return k.value
}

// ボタンの種類が存在しているか確認します
func (k Kind) IsEmpty() bool {
	return k.value == ""
}

// ボタンの種類を検証します
func (k Kind) validate() error {
	switch k.value {
	case ButtonKindToOpen:
	case ButtonKindToResult:
	default:
		return errors.NewError("ボタンの種類が不正です", nil)
	}

	return nil
}

// ボタンの種類をJSONに変換します
func (k Kind) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: k.value,
	}

	return json.Marshal(data)
}

// ボタンの種類をJSONから変換します
func (k *Kind) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	k.value = data.Value

	return nil
}
