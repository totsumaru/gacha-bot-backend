package user_data

import (
	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/user_data"
	gatewayUserData "github.com/totsumaru/gacha-bot-backend/gateway/user_data"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// Upsertのリクエストです
type UpsertReq struct {
	ServerID string
	UserID   string
	Point    int
}

// ユーザーデータをUpsertします
func UpsertUserData(tx *gorm.DB, req UpsertReq) (user_data.UserData, error) {
	serverID, err := domain.NewDiscordID(req.ServerID)
	if err != nil {
		return user_data.UserData{}, errors.NewError("サーバーIDの生成に失敗しました", err)
	}

	userID, err := domain.NewDiscordID(req.UserID)
	if err != nil {
		return user_data.UserData{}, errors.NewError("ユーザーIDの生成に失敗しました", err)
	}

	point, err := user_data.NewPoint(req.Point)
	if err != nil {
		return user_data.UserData{}, errors.NewError("ポイントの生成に失敗しました", err)
	}

	ud, err := user_data.NewUserData(serverID, userID, point)
	if err != nil {
		return user_data.UserData{}, errors.NewError("ユーザーデータの生成に失敗しました", err)
	}

	gw, err := gatewayUserData.NewGateway(tx)
	if err != nil {
		return user_data.UserData{}, errors.NewError("Gatewayの生成に失敗しました", err)
	}

	if err = gw.Upsert(ud); err != nil {
		return user_data.UserData{}, errors.NewError("ユーザーデータの更新に失敗しました", err)
	}

	return ud, nil
}

// ポイントを追加します
func AddPoint(tx *gorm.DB, serverID, userID string, addPoint int) error {
	sID, err := domain.NewDiscordID(serverID)
	if err != nil {
		return errors.NewError("サーバーIDの生成に失敗しました", err)
	}

	uID, err := domain.NewDiscordID(userID)
	if err != nil {
		return errors.NewError("ユーザーIDの生成に失敗しました", err)
	}

	id, err := user_data.NewID(sID, uID)
	if err != nil {
		return errors.NewError("IDの生成に失敗しました", err)
	}

	gw, err := gatewayUserData.NewGateway(tx)
	if err != nil {
		return errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	currentPoint := 0

	ud, err := gw.FindByIDForUpdate(id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			// レコードが存在しない場合
			currentPoint = 0
		} else {
			return errors.NewError("IDでユーザーデータを取得できません", err)
		}
	} else {
		currentPoint = ud.Point().Int()
	}

	newPoint := currentPoint + addPoint

	p, err := user_data.NewPoint(newPoint)
	if err != nil {
		return errors.NewError("ポイントを作成できません", err)
	}

	newUserData, err := user_data.NewUserData(sID, uID, p)
	if err != nil {
		return errors.NewError("ユーザーデータを作成できません", err)
	}

	if err = gw.Upsert(newUserData); err != nil {
		return errors.NewError("ユーザーデータの更新に失敗しました", err)
	}

	return nil
}

// サーバーIDとユーザーIDでユーザーデータを取得します
func FindByServerIDAndUserID(tx *gorm.DB, serverID, userID string) (user_data.UserData, error) {
	sID, err := domain.NewDiscordID(serverID)
	if err != nil {
		return user_data.UserData{}, errors.NewError("サーバーIDの生成に失敗しました", err)
	}

	uID, err := domain.NewDiscordID(userID)
	if err != nil {
		return user_data.UserData{}, errors.NewError("ユーザーIDの生成に失敗しました", err)
	}

	id, err := user_data.NewID(sID, uID)
	if err != nil {
		return user_data.UserData{}, errors.NewError("IDの生成に失敗しました", err)
	}

	gw, err := gatewayUserData.NewGateway(tx)
	if err != nil {
		return user_data.UserData{}, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	s, err := gw.FindByID(id)
	if err != nil {
		return user_data.UserData{}, errors.NewError("IDでサーバーを取得できません", err)
	}

	return s, nil
}

// FOR UPDATEでサーバーIDとユーザーIDでユーザーデータを取得します
func FindByServerIDAndUserIDForUpdate(tx *gorm.DB, serverID, userID string) (user_data.UserData, error) {
	sID, err := domain.NewDiscordID(serverID)
	if err != nil {
		return user_data.UserData{}, errors.NewError("サーバーIDの生成に失敗しました", err)
	}

	uID, err := domain.NewDiscordID(userID)
	if err != nil {
		return user_data.UserData{}, errors.NewError("ユーザーIDの生成に失敗しました", err)
	}

	id, err := user_data.NewID(sID, uID)
	if err != nil {
		return user_data.UserData{}, errors.NewError("IDの生成に失敗しました", err)
	}

	gw, err := gatewayUserData.NewGateway(tx)
	if err != nil {
		return user_data.UserData{}, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	s, err := gw.FindByIDForUpdate(id)
	if err != nil {
		return user_data.UserData{}, errors.NewError("IDでサーバーを取得できません", err)
	}

	return s, nil
}
