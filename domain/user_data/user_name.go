package user_data

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ユーザー名です
type UserName struct {
	value string
}

// ユーザー名を作成します
func NewUserName(value string) (UserName, error) {
	u := UserName{
		value: value,
	}

	if err := u.validate(); err != nil {
		return UserName{}, errors.NewError("ユーザー名の生成に失敗しました", err)
	}

	return u, nil
}

// 値を取得します
func (u UserName) String() string {
	return u.value
}

// 値が空かどうかを返します
func (u UserName) IsEmpty() bool {
	return u.value == ""
}

// 値を検証します
func (u UserName) validate() error {
	return nil
}

// 値をJSONに変換します
func (u UserName) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: u.value,
	}

	return json.Marshal(data)
}

// JSONから値を復元します
func (u *UserName) UnmarshalJSON(b []byte) error {
	data := struct {
		Value string `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("UserNameの復元に失敗しました", err)
	}

	u.value = data.Value

	return nil
}
