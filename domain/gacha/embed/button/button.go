package button

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ボタンコンポーネントです
type Button struct {
	kind  Kind
	label Label
	style Style
}

// ボタンを作成します
func NewButton(kind Kind, label Label, style Style) (Button, error) {
	b := Button{
		kind:  kind,
		label: label,
		style: style,
	}

	if err := b.validate(); err != nil {
		return b, errors.NewError("検証に失敗しました", err)
	}

	return b, nil
}

// IDを返します
func (b Button) Kind() Kind {
	return b.kind
}

// ボタンのラベルを返します
func (b Button) Label() Label {
	return b.label
}

// ボタンのスタイルを返します
func (b Button) Style() Style {
	return b.style
}

// 検証します
func (b Button) validate() error {
	return nil
}

// MarshalJSON は Button 構造体を JSON に変換します。
func (b Button) MarshalJSON() ([]byte, error) {
	bb, err := json.Marshal(struct {
		Kind  Kind  `json:"kind"`
		Label Label `json:"label"`
		Style Style `json:"style"`
	}{
		Kind:  b.kind,
		Label: b.label,
		Style: b.style,
	})

	if err != nil {
		return nil, errors.NewError("JSONに変換できませんでした", err)
	}

	return bb, nil
}

// UnmarshalJSON は JSON から Button 構造体を復元します。
func (b *Button) UnmarshalJSON(bytes []byte) error {
	data := struct {
		Kind  Kind  `json:"kind"`
		Label Label `json:"label"`
		Style Style `json:"style"`
	}{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return errors.NewError("JSONからボタンの復元に失敗しました", err)
	}

	b.kind = data.Kind
	b.label = data.Label
	b.style = data.Style

	return nil
}
