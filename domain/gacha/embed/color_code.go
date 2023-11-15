package embed

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/color"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// カラーコードです
type ColorCode struct {
	value int
}

// カラーコードを生成します
func NewColorCode(value int) (ColorCode, error) {
	c := ColorCode{value: value}

	if err := c.validate(); err != nil {
		return ColorCode{}, errors.NewError("カラーコードが不正です", err)
	}

	return c, nil
}

// カラーコードを返します
func (c ColorCode) Int() int {
	return c.value
}

// カラーコードが存在しているか確認します
func (c ColorCode) IsZero() bool {
	return c.value == 0
}

// カラーコードを検証します
func (c ColorCode) validate() error {
	switch c.value {
	case color.Blue:
	case color.Red:
	case color.Orange:
	case color.Green:
	case color.Pink:
	case color.Black:
	case color.Yellow:
	case color.Cyan:
	case color.DarkGray:
	default:
		return errors.NewError("カラーコードが不正です", nil)
	}

	return nil
}

// カラーコードをJSONに変換します
func (c ColorCode) MarshalJSON() ([]byte, error) {
	data := struct {
		Value int `json:"value"`
	}{
		Value: c.value,
	}

	return json.Marshal(data)
}

// カラーコードをJSONから変換します
func (c *ColorCode) UnmarshalJSON(b []byte) error {
	data := struct {
		Value int `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	c.value = data.Value

	return nil
}
