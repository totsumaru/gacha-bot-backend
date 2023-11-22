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
	id       domain.UUID
	serverID domain.DiscordID
	panel    embed.Embed
	open     embed.Embed
	result   []result.Result
	role     []Role
}

// ガチャを生成します
func NewGacha(
	serverID domain.DiscordID,
	panel embed.Embed,
	open embed.Embed,
	result []result.Result,
	role []Role,
) (Gacha, error) {
	id, err := domain.NewUUID()
	if err != nil {
		return Gacha{}, errors.NewError("IDの生成に失敗しました", err)
	}

	g := Gacha{
		id:       id,
		serverID: serverID,
		panel:    panel,
		open:     open,
		result:   result,
		role:     role,
	}

	if err = g.validate(); err != nil {
		return Gacha{}, errors.NewError("ガチャの生成に失敗しました", err)
	}

	return g, nil
}

// ガチャを復元します
func RestoreGacha(
	id domain.UUID,
	serverID domain.DiscordID,
	panel embed.Embed,
	open embed.Embed,
	result []result.Result,
	role []Role,
) (Gacha, error) {
	g := Gacha{
		id:       id,
		serverID: serverID,
		panel:    panel,
		open:     open,
		result:   result,
		role:     role,
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

// サーバーIDを返します
func (g Gacha) ServerID() domain.DiscordID {
	return g.serverID
}

// パネルを返します
func (g Gacha) Panel() embed.Embed {
	return g.panel
}

// オープンを返します
func (g Gacha) Open() embed.Embed {
	return g.open
}

// 結果を返します
func (g Gacha) Result() []result.Result {
	return g.result
}

// ロールを返します
func (g Gacha) Role() []Role {
	return g.role
}

// ガチャを検証します
func (g Gacha) validate() error {
	// 確率の合計が100%か確認します
	sum := 0
	for _, r := range g.result {
		sum += r.Probability().Int()
	}

	if sum != 100 {
		return errors.NewError("確率の合計が100％ではありません")
	}

	// ロールIDが重複していないか確認します
	ids := map[string]bool{}
	for _, r := range g.role {
		if _, ok := ids[r.ID().String()]; ok {
			return errors.NewError("ロールIDが重複しています")
		} else {
			ids[r.ID().String()] = true
		}
	}

	return nil
}

// ガチャをJSONに変換します
func (g Gacha) MarshalJSON() ([]byte, error) {
	data := struct {
		ID       domain.UUID      `json:"id"`
		ServerID domain.DiscordID `json:"server_id"`
		Panel    embed.Embed      `json:"panel"`
		Open     embed.Embed      `json:"open"`
		Result   []result.Result  `json:"result"`
		Role     []Role           `json:"role"`
	}{
		ID:       g.id,
		ServerID: g.serverID,
		Panel:    g.panel,
		Open:     g.open,
		Result:   g.result,
		Role:     g.role,
	}

	return json.Marshal(data)
}

// ガチャをJSONから復元します
func (g *Gacha) UnmarshalJSON(b []byte) error {
	data := struct {
		ID       domain.UUID      `json:"id"`
		ServerID domain.DiscordID `json:"server_id"`
		Panel    embed.Embed      `json:"panel"`
		Open     embed.Embed      `json:"open"`
		Result   []result.Result  `json:"result"`
		Role     []Role           `json:"role"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	g.id = data.ID
	g.serverID = data.ServerID
	g.panel = data.Panel
	g.open = data.Open
	g.result = data.Result
	g.role = data.Role

	return nil
}
