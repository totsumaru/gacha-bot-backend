package user_data

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ユーザーデータです
type UserData struct {
	id       ID
	serverID domain.DiscordID
	userID   domain.DiscordID
	point    Point
}

// ユーザーデータを作成します
func NewUserData(
	serverID domain.DiscordID,
	userID domain.DiscordID,
	point Point,
) (UserData, error) {
	id, err := NewID(serverID, userID)
	if err != nil {
		return UserData{}, errors.NewError("ユーザーデータの生成に失敗しました", err)
	}

	p := UserData{
		id:       id,
		serverID: serverID,
		userID:   userID,
		point:    point,
	}

	if err = p.validate(); err != nil {
		return UserData{}, errors.NewError("ユーザーデータの生成に失敗しました", err)
	}

	return p, nil
}

// ポイントを更新します
func (p UserData) UpdatePoint(point Point) error {
	p.point = point
	if err := p.validate(); err != nil {
		return errors.NewError("ユーザーデータの更新に失敗しました", err)
	}

	return nil
}

// IDを返します
func (p UserData) ID() ID {
	return p.id
}

// サーバーIDを返します
func (p UserData) ServerID() domain.DiscordID {
	return p.serverID
}

// ユーザーIDを返します
func (p UserData) UserID() domain.DiscordID {
	return p.userID
}

// ポイントを返します
func (p UserData) Point() Point {
	return p.point
}

// 検証します
func (p UserData) validate() error {
	return nil
}

// ユーザーデータをJSONに変換します
func (p UserData) MarshalJSON() ([]byte, error) {
	data := struct {
		ID       ID               `json:"id"`
		ServerID domain.DiscordID `json:"server_id"`
		UserID   domain.DiscordID `json:"user_id"`
		Point    Point            `json:"point"`
	}{
		ID:       p.id,
		ServerID: p.serverID,
		UserID:   p.userID,
		Point:    p.point,
	}

	return json.Marshal(data)
}

// JSONからユーザーデータを作成します
func (p *UserData) UnmarshalJSON(b []byte) error {
	data := struct {
		ID       ID               `json:"id"`
		ServerID domain.DiscordID `json:"server_id"`
		UserID   domain.DiscordID `json:"user_id"`
		Point    Point            `json:"point"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("ユーザーデータのJSONのパースに失敗しました", err)
	}

	p.id = data.ID
	p.serverID = data.ServerID
	p.userID = data.UserID
	p.point = data.Point

	return nil
}
