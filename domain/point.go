package domain

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ポイントです
type Point struct {
	value int
}

// ポイントを生成します
func NewPoint(value int) (Point, error) {
	p := Point{value: value}

	if err := p.validate(); err != nil {
		return Point{}, errors.NewError("ポイントが不正です", err)
	}

	return p, nil
}

// ポイントを返します
func (p Point) Int() int {
	return p.value
}

// ポイントが存在しているか確認します
func (p Point) IsZero() bool {
	return p.value == 0
}

// ポイントを検証します
func (p Point) validate() error {
	if p.value < 0 {
		return errors.NewError("ポイントがマイナスの値です")
	}

	return nil
}

// ポイントをJSONに変換します
func (p Point) MarshalJSON() ([]byte, error) {
	data := struct {
		Value int `json:"value"`
	}{
		Value: p.value,
	}

	return json.Marshal(data)
}

// ポイントをJSONから変換します
func (p *Point) UnmarshalJSON(b []byte) error {
	data := struct {
		Value int `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	p.value = data.Value

	return nil
}
