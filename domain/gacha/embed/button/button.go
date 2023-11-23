package button

import (
	"encoding/json"
	"fmt"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ボタンコンポーネントです
type Button struct {
	kind     Kind
	label    Label
	style    Style
	url      domain.URL
	isHidden bool // ボタンの表示/非表示
}

// ボタンを作成します
func NewButton(
	kind Kind,
	label Label,
	style Style,
	url domain.URL,
	isHidden bool,
) (Button, error) {
	b := Button{
		kind:     kind,
		label:    label,
		style:    style,
		url:      url,
		isHidden: isHidden,
	}

	// ボタンのスタイルがLink以外の場合は、URLを空にする
	if b.Style().String() != ButtonStyleLink {
		b.url = domain.URL{}
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

// URLを返します
func (b Button) URL() domain.URL {
	return b.url
}

// ボタンが非表示かどうかを返します
func (b Button) IsHidden() bool {
	return b.isHidden
}

// 検証します
func (b Button) validate() error {
	fmt.Println("isHidden", b.isHidden, ": Style: ", b.Style().String())
	switch b.Style().String() {
	case ButtonStyleLink:
		if !b.isHidden {
			if b.Label().IsEmpty() {
				return errors.NewError("Linkボタンにはラベルが必須です")
			}
			if b.URL().IsEmpty() {
				return errors.NewError("LinkボタンにはURLが必須です")
			}
		}
	default:
		if b.Label().IsEmpty() {
			return errors.NewError("ラベルが必須です")
		}
	}

	return nil
}

// MarshalJSON は Button 構造体を JSON に変換します。
func (b Button) MarshalJSON() ([]byte, error) {
	bb, err := json.Marshal(struct {
		Kind     Kind       `json:"kind"`
		Label    Label      `json:"label"`
		Style    Style      `json:"style"`
		URL      domain.URL `json:"url"`
		IsHidden bool       `json:"is_hidden"`
	}{
		Kind:     b.kind,
		Label:    b.label,
		Style:    b.style,
		URL:      b.url,
		IsHidden: b.isHidden,
	})

	if err != nil {
		return nil, errors.NewError("JSONに変換できませんでした", err)
	}

	return bb, nil
}

// UnmarshalJSON は JSON から Button 構造体を復元します。
func (b *Button) UnmarshalJSON(bytes []byte) error {
	data := struct {
		Kind     Kind       `json:"kind"`
		Label    Label      `json:"label"`
		Style    Style      `json:"style"`
		URL      domain.URL `json:"url"`
		IsHidden bool       `json:"is_hidden"`
	}{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return errors.NewError("JSONからボタンの復元に失敗しました", err)
	}

	b.kind = data.Kind
	b.label = data.Label
	b.style = data.Style
	b.url = data.URL
	b.isHidden = data.IsHidden

	return nil
}
