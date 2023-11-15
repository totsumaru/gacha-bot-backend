package domain

import (
	"encoding/json"
	"net/url"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// URLです
type URL struct {
	value string
}

// URLを作成します
func NewURL(value string) (URL, error) {
	u := URL{value: value}

	if err := u.validate(); err != nil {
		return u, errors.NewError("検証に失敗しました", err)
	}

	return u, nil
}

// URLを返します
func (u URL) String() string {
	return u.value
}

// URLが存在しているか確認します
func (u URL) IsEmpty() bool {
	return u.value == ""
}

// URLを検証します
func (u URL) validate() error {
	if u.IsEmpty() {
		return nil
	}

	if !isValidURL(u.value) {
		return errors.NewError("URLが不正です")
	}

	return nil
}

// URLの形式になっていることを確認します
func isValidURL(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

// URLをJSONに変換します
func (u URL) MarshalJSON() ([]byte, error) {
	data := struct {
		URL string `json:"kind"`
	}{
		URL: u.value,
	}

	return json.Marshal(data)
}

// URLをJSONから復元します
func (u *URL) UnmarshalJSON(b []byte) error {
	data := struct {
		URL string `json:"kind"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	u.value = data.URL

	return nil
}
