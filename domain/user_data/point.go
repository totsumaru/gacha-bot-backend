package user_data

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ポイントの合計です
type Point struct {
	value int
}

// ポイントの合計を作成します
func NewPoint(value int) (Point, error) {
	p := Point{
		value: value,
	}

	if err := p.validate(); err != nil {
		return Point{}, errors.NewError("ポイントの合計の生成に失敗しました", err)
	}

	return p, nil
}

// ポイントの合計を返します
func (p Point) String() int {
	return p.value
}

// ポイントの合計を検証します
func (p Point) validate() error {
	if p.value < 0 {
		return errors.NewError("ポイントの合計が不正です")
	}

	return nil
}

// ポイントの合計をJSONに変換します
func (p Point) MarshalJSON() ([]byte, error) {
	data := struct {
		Value int `json:"value"`
	}{
		Value: p.value,
	}

	return json.Marshal(data)
}

// JSONからポイントの合計を作成します
func (p *Point) UnmarshalJSON(b []byte) error {
	data := struct {
		Value int `json:"value"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("ポイントの合計のJSONのパースに失敗しました", err)
	}

	p.value = data.Value

	return nil
}
