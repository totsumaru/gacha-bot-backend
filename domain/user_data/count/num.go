package count

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// 1日の回数です
type Num struct {
	value int
}

// 回数を生成します
func NewNum(value int) (Num, error) {
	n := Num{value: value}

	if err := n.validate(); err != nil {
		return Num{}, errors.NewError("回数が不正です", err)
	}

	return n, nil
}

// 回数を返します
func (n Num) Int() int {
	return n.value
}

// 回数が存在しているか確認します
func (n Num) IsZero() bool {
	return n.value == 0
}

// 回数を検証します
func (n Num) validate() error {
	return nil
}

// 回数をJSONに変換します
func (n Num) MarshalJSON() ([]byte, error) {
	data := struct {
		Value int `json:"value"`
	}{
		Value: n.value,
	}

	return json.Marshal(data)
}

// 回数をJSONから変換します
func (n *Num) UnmarshalJSON(b []byte) error {
	data := struct {
		Value int `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	n.value = data.Value

	return nil
}
