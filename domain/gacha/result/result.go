package result

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// 結果の埋め込みです
type Result struct {
	embed       embed.Embed
	point       domain.Point
	probability Probability
}

// 結果を生成します
func NewResult(
	embed embed.Embed,
	point domain.Point,
	probability Probability,
) (Result, error) {
	r := Result{
		embed:       embed,
		point:       point,
		probability: probability,
	}

	if err := r.validate(); err != nil {
		return Result{}, errors.NewError("結果の生成に失敗しました", err)
	}

	return r, nil
}

// 埋め込みを返します
func (r Result) Embed() embed.Embed {
	return r.embed
}

// ポイントを返します
func (r Result) Point() domain.Point {
	return r.point
}

// 確率を返します
func (r Result) Probability() Probability {
	return r.probability
}

// 結果を検証します
func (r Result) validate() error {
	return nil
}

// JSONに変換します
func (r Result) MarshalJSON() ([]byte, error) {
	data := struct {
		Embed       embed.Embed  `json:"embed"`
		Point       domain.Point `json:"point"`
		Probability Probability  `json:"probability"`
	}{
		Embed:       r.embed,
		Point:       r.point,
		Probability: r.probability,
	}

	return json.Marshal(data)
}

// JSONから変換します
func (r *Result) UnmarshalJSON(b []byte) error {
	data := struct {
		Embed       embed.Embed  `json:"embed"`
		Point       domain.Point `json:"point"`
		Probability Probability  `json:"probability"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	r.embed = data.Embed
	r.point = data.Point
	r.probability = data.Probability

	return nil
}
