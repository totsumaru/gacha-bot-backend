package stripe

import (
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
func (i SubscriptionID) Value() string {
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
