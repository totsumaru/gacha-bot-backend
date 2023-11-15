package button

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

const (
	ButtonStylePrimary   = "PRIMARY"
	ButtonStyleSecondary = "SECONDARY"
	ButtonStyleSuccess   = "SUCCESS"
	ButtonStyleDanger    = "DANGER"
	ButtonStyleLink      = "LINK"
)

// ボタンのスタイルです
type Style struct {
	value string
}

// ボタンのスタイルを作成します
func NewStyle(value string) (Style, error) {
	s := Style{value: value}

	if err := s.validate(); err != nil {
		return s, err
	}

	return s, nil
}

// ボタンのスタイルを返します
func (s Style) String() string {
	return s.value
}

// ボタンのスタイルが存在しているか確認します
func (s Style) IsEmpty() bool {
	return s.value == ""
}

// ボタンのスタイルを検証します
func (s Style) validate() error {
	switch s.value {
	case ButtonStylePrimary:
	case ButtonStyleSecondary:
	case ButtonStyleSuccess:
	case ButtonStyleDanger:
	case ButtonStyleLink:
	default:
		return errors.NewError("ボタンのスタイルが不正です")
	}

	return nil
}

// ボタンのスタイルをJSONに変換します
func (s Style) MarshalJSON() ([]byte, error) {
	data := struct {
		Style string `json:"style"`
	}{
		Style: s.value,
	}

	return json.Marshal(data)
}

// JSONからボタンのスタイルを復元します
func (s *Style) UnmarshalJSON(b []byte) error {
	data := struct {
		Style string `json:"style"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONからボタンのスタイルの復元に失敗しました", err)
	}

	s.value = data.Style

	return nil
}
