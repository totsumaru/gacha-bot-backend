package stripe

import (
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
func (i CustomerID) Value() string {
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
