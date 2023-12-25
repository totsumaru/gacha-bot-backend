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
	userName UserName
	iconURL  IconURL
}

// ユーザーデータを作成します
func NewUserData(
	serverID domain.DiscordID,
	userID domain.DiscordID,
	point domain.Point,
	count count.Count,
	userName UserName,
	iconURL IconURL,
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
		userName: userName,
		iconURL:  iconURL,
	}

	if err = p.validate(); err != nil {
		return UserData{}, errors.NewError("ユーザーデータの生成に失敗しました", err)
	}

	return p, nil
}

// ポイントを更新します
func (u *UserData) UpdatePoint(point domain.Point) error {
	u.point = point
	if err := u.validate(); err != nil {
		return errors.NewError("ポイントの更新に失敗しました", err)
	}

	return nil
}

// 1日の上限回数を更新します
func (u *UserData) UpdateCount(count count.Count) error {
	u.count = count
	if err := u.validate(); err != nil {
		return errors.NewError("1日の上限回数の更新に失敗しました", err)
	}

	return nil
}

// ユーザー名を更新します
func (u *UserData) UpdateUserName(userName UserName) error {
	u.userName = userName
	if err := u.validate(); err != nil {
		return errors.NewError("ユーザー名の更新に失敗しました", err)
	}

	return nil
}

// アイコンURLを更新します
func (u *UserData) UpdateIconURL(iconURL IconURL) error {
	u.iconURL = iconURL
	if err := u.validate(); err != nil {
		return errors.NewError("アイコンURLの更新に失敗しました", err)
	}

	return nil
}

// IDを返します
func (u UserData) ID() ID {
	return u.id
}

// サーバーIDを返します
func (u UserData) ServerID() domain.DiscordID {
	return u.serverID
}

// ユーザーIDを返します
func (u UserData) UserID() domain.DiscordID {
	return u.userID
}

// ポイントを返します
func (u UserData) Point() domain.Point {
	return u.point
}

// 1日の上限回数を返します
func (u UserData) Count() count.Count {
	return u.count
}

// ユーザー名を返します
func (u UserData) UserName() UserName {
	return u.userName
}

// アイコンURLを返します
func (u UserData) IconURL() IconURL {
	return u.iconURL
}

// 検証します
func (u UserData) validate() error {
	return nil
}

// ユーザーデータをJSONに変換します
func (u UserData) MarshalJSON() ([]byte, error) {
	data := struct {
		ID       ID               `json:"id"`
		ServerID domain.DiscordID `json:"server_id"`
		UserID   domain.DiscordID `json:"user_id"`
		Point    domain.Point     `json:"point"`
		Count    count.Count      `json:"count"`
		UserName UserName         `json:"user_name"`
		IconURL  IconURL          `json:"icon_url"`
	}{
		ID:       u.id,
		ServerID: u.serverID,
		UserID:   u.userID,
		Point:    u.point,
		Count:    u.count,
		UserName: u.userName,
		IconURL:  u.iconURL,
	}

	return json.Marshal(data)
}

// JSONからユーザーデータを作成します
func (u *UserData) UnmarshalJSON(b []byte) error {
	data := struct {
		ID       ID               `json:"id"`
		ServerID domain.DiscordID `json:"server_id"`
		UserID   domain.DiscordID `json:"user_id"`
		Point    domain.Point     `json:"point"`
		Count    count.Count      `json:"count"`
		UserName UserName         `json:"user_name"`
		IconURL  IconURL          `json:"icon_url"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("ユーザーデータのJSONのパースに失敗しました", err)
	}

	u.id = data.ID
	u.serverID = data.ServerID
	u.userID = data.UserID
	u.point = data.Point
	u.count = data.Count
	u.userName = data.UserName
	u.iconURL = data.IconURL

	return nil
}
