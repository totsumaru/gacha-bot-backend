package user_data

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/user_data/count"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ユーザーデータです
type UserData struct {
	id       ID
	serverID domain.DiscordID
	userID   domain.DiscordID
	point    domain.Point
	count    count.Count
}

// ユーザーデータを作成します
func NewUserData(
	serverID domain.DiscordID,
	userID domain.DiscordID,
	point domain.Point,
	count count.Count,
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
		count:    count,
	}

	if err = p.validate(); err != nil {
		return UserData{}, errors.NewError("ユーザーデータの生成に失敗しました", err)
	}

	return p, nil
}

// ポイントを更新します
func (p UserData) UpdatePoint(point domain.Point) error {
	p.point = point
	if err := p.validate(); err != nil {
		return errors.NewError("ポイントの更新に失敗しました", err)
	}

	return nil
}

// 1日の上限回数を更新します
func (p UserData) UpdateCount(count count.Count) error {
	p.count = count
	if err := p.validate(); err != nil {
		return errors.NewError("1日の上限回数の更新に失敗しました", err)
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
func (p UserData) Point() domain.Point {
	return p.point
}

// 1日の上限回数を返します
func (p UserData) Count() count.Count {
	return p.count
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
		Point    domain.Point     `json:"point"`
		Count    count.Count      `json:"count"`
	}{
		ID:       p.id,
		ServerID: p.serverID,
		UserID:   p.userID,
		Point:    p.point,
		Count:    p.count,
	}

	return json.Marshal(data)
}

// JSONからユーザーデータを作成します
func (p *UserData) UnmarshalJSON(b []byte) error {
	data := struct {
		ID       ID               `json:"id"`
		ServerID domain.DiscordID `json:"server_id"`
		UserID   domain.DiscordID `json:"user_id"`
		Point    domain.Point     `json:"point"`
		Count    count.Count      `json:"count"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("ユーザーデータのJSONのパースに失敗しました", err)
	}

	p.id = data.ID
	p.serverID = data.ServerID
	p.userID = data.UserID
	p.point = data.Point
	p.count = data.Count

	return nil
}
