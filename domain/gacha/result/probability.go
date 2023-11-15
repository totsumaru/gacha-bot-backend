package result

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// 確率です
//
// 確率は0以上100以下の整数です。
type Probability struct {
	value int
}

// 確率を生成します
func NewProbability(value int) (Probability, error) {
	p := Probability{value: value}

	if err := p.validate(); err != nil {
		return Probability{}, errors.NewError("確率が不正です", err)
	}

	return p, nil
}

// 確率を返します
func (p Probability) Int() int {
	return p.value
}

// 確率が存在しているか確認します
func (p Probability) IsZero() bool {
	return p.value == 0
}

// 確率を検証します
func (p Probability) validate() error {
	if p.value < 0 || p.value > 100 {
		return errors.NewError("確率が不正です", nil)
	}

	return nil
}

// 確率をJSONに変換します
func (p Probability) MarshalJSON() ([]byte, error) {
	data := struct {
		Value int `json:"value"`
	}{
		Value: p.value,
	}

	return json.Marshal(data)
}

// 確率をJSONから変換します
func (p *Probability) UnmarshalJSON(b []byte) error {
	data := struct {
		Value int `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	p.value = data.Value

	return nil
}
