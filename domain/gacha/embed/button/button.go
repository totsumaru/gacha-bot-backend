package button

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ボタンコンポーネントです
type Button struct {
	id    domain.UUID
	label Label
	style Style
}

// ボタンを作成します
func NewButton(id domain.UUID, label Label, style Style) (Button, error) {
	b := Button{
		id:    id,
		label: label,
		style: style,
	}

	if err := b.validate(); err != nil {
		return b, errors.NewError("検証に失敗しました", err)
	}

	return b, nil
}

// IDを返します
func (b Button) ID() domain.UUID {
	return b.id
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
		ID    domain.UUID `json:"id"`
		Label Label       `json:"label"`
		Style Style       `json:"style"`
	}{
		ID:    b.ID(),
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
		ID    domain.UUID `json:"id"`
		Label Label       `json:"label"`
		Style Style       `json:"style"`
	}{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return errors.NewError("JSONからボタンの復元に失敗しました", err)
	}

	b.id = data.ID
	b.label = data.Label
	b.style = data.Style

	return nil
}
