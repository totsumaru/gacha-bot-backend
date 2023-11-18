package stripe

import (
	"encoding/json"
	"strings"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// カスタマーIDです
type CustomerID struct {
	value string
}

// カスタマーIDを作成します
func NewCustomerID(v string) (CustomerID, error) {
	id := CustomerID{}
	id.value = v

	if err := id.validate(); err != nil {
		return id, errors.NewError("検証に失敗しました", err)
	}

	return id, nil
}

// カスタマーIDを取得します
func (i CustomerID) String() string {
	return i.value
}

// 検証をします
func (i CustomerID) validate() error {
	if i.value == "" {
		return nil
	}

	if !strings.HasPrefix(i.value, "cus_") {
		return errors.NewError("指定した文字列から始まっていません")
	}

	return nil
}

// JSONに変換します
func (i CustomerID) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: i.value,
	}

	return json.Marshal(data)
}

// JSONから復元します
func (i *CustomerID) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONから復元できません", err)
	}

	i.value = data.Value

	if err := i.validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}
