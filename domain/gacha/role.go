package gacha

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
)

// ロールです
type Role struct {
	id    domain.DiscordID
	point domain.Point
}

// ロールを生成します
func NewRole(id domain.DiscordID, point domain.Point) (Role, error) {
	r := Role{id: id, point: point}

	if err := r.validate(); err != nil {
		return Role{}, err
	}

	return r, nil
}

// ロールIDを返します
func (r Role) ID() domain.DiscordID {
	return r.id
}

// ポイントを返します
func (r Role) Point() domain.Point {
	return r.point
}

// ロールを検証します
func (r Role) validate() error {
	return nil
}

// ロールをJSONに変換します
func (r Role) MarshalJSON() ([]byte, error) {
	data := struct {
		ID    domain.DiscordID `json:"id"`
		Point domain.Point     `json:"point"`
	}{
		ID:    r.id,
		Point: r.point,
	}

	return json.Marshal(data)
}

// ロールをJSONから変換します
func (r *Role) UnmarshalJSON(b []byte) error {
	data := struct {
		ID    domain.DiscordID `json:"id"`
		Point domain.Point     `json:"point"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	r.id = data.ID
	r.point = data.Point

	return nil
}
