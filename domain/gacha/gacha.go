package gacha

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/result"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ガチャのドメインモデルです
type Gacha struct {
	id     domain.UUID
	panel  embed.Embed
	open   embed.Embed
	result []result.Result
}

// ガチャを生成します
func NewGacha(
	panel embed.Embed,
	open embed.Embed,
	result []result.Result,
) (Gacha, error) {
	id, err := domain.NewUUID()
	if err != nil {
		return Gacha{}, errors.NewError("IDの生成に失敗しました", err)
	}

	g := Gacha{
		id:     id,
		panel:  panel,
		open:   open,
		result: result,
	}

	if err = g.validate(); err != nil {
		return Gacha{}, errors.NewError("ガチャの生成に失敗しました", err)
	}

	return g, nil
}

// ガチャを復元します
func RestoreGacha(
	id domain.UUID,
	panel embed.Embed,
	open embed.Embed,
	result []result.Result,
) (Gacha, error) {
	g := Gacha{
		id:     id,
		panel:  panel,
		open:   open,
		result: result,
	}

	if err := g.validate(); err != nil {
		return Gacha{}, errors.NewError("ガチャの生成に失敗しました", err)
	}

	return g, nil
}

// IDを返します
func (g Gacha) ID() domain.UUID {
	return g.id
}

// パネルを返します
func (g Gacha) Panel() embed.Embed {
	return g.panel
}

// オープニングを返します
func (g Gacha) Open() embed.Embed {
	return g.open
}

// 結果を返します
func (g Gacha) Result() []result.Result {
	return g.result
}

// ガチャを検証します
func (g Gacha) validate() error {
	// 確率の合計が100%か確認します
	sum := 0
	for _, r := range g.result {
		sum += r.Probability().Int()
	}

	if sum != 100 {
		return errors.NewError("確率の合計が100%ではありません", nil)
	}

	return nil
}

// ガチャをJSONに変換します
func (g Gacha) MarshalJSON() ([]byte, error) {
	data := struct {
		ID     domain.UUID     `json:"id"`
		Panel  embed.Embed     `json:"panel"`
		Open   embed.Embed     `json:"open"`
		Result []result.Result `json:"result"`
	}{
		ID:     g.id,
		Panel:  g.panel,
		Open:   g.open,
		Result: g.result,
	}

	return json.Marshal(data)
}

// ガチャをJSONから復元します
func (g *Gacha) UnmarshalJSON(b []byte) error {
	data := struct {
		ID     domain.UUID     `json:"id"`
		Panel  embed.Embed     `json:"panel"`
		Open   embed.Embed     `json:"open"`
		Result []result.Result `json:"result"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	g.id = data.ID
	g.panel = data.Panel
	g.open = data.Open
	g.result = data.Result

	return nil
}
