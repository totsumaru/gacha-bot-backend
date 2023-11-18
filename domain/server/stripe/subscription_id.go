package stripe

import (
	"encoding/json"
	"strings"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// サブスクリプションIDです
type SubscriptionID struct {
	value string
}

// サブスクリプションIDを作成します
func NewSubscriptionID(v string) (SubscriptionID, error) {
	id := SubscriptionID{}
	id.value = v

	if err := id.validate(); err != nil {
		return id, errors.NewError("検証に失敗しました", err)
	}

	return id, nil
}

// サブスクリプションIDを取得します
func (i SubscriptionID) String() string {
	return i.value
}

// 検証をします
func (i SubscriptionID) validate() error {
	if i.value == "" {
		return nil
	}

	if !strings.HasPrefix(i.value, "sub_") {
		return errors.NewError("指定した文字列から始まっていません")
	}

	return nil
}

// JSONに変換します
func (i SubscriptionID) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: i.value,
	}

	return json.Marshal(data)
}

// JSONから復元します
func (i *SubscriptionID) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONから復元できませんでした", err)
	}

	i.value = data.Value

	if err := i.validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}
