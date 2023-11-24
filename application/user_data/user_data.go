package user_data

import (
	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/user_data"
	"github.com/totsumaru/gacha-bot-backend/domain/user_data/count"
	gatewayUserData "github.com/totsumaru/gacha-bot-backend/gateway/user_data"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"github.com/totsumaru/gacha-bot-backend/lib/now"
	"gorm.io/gorm"
)

// ポイントを追加します
//
// 追加後のポイント(最新のポイント)を返します。
func AddPoint(tx *gorm.DB, serverID, userID string, addPoint int) (int, error) {
	sID, err := domain.NewDiscordID(serverID)
	if err != nil {
		return 0, errors.NewError("サーバーIDの生成に失敗しました", err)
	}

	uID, err := domain.NewDiscordID(userID)
	if err != nil {
		return 0, errors.NewError("ユーザーIDの生成に失敗しました", err)
	}

	id, err := user_data.NewID(sID, uID)
	if err != nil {
		return 0, errors.NewError("IDの生成に失敗しました", err)
	}

	gw, err := gatewayUserData.NewGateway(tx)
	if err != nil {
		return 0, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	currentPoint := 0
	lim := count.Count{}

	ud, err := gw.FindByIDForUpdate(id)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			return 0, errors.NewError("IDでユーザーデータを取得できません", err)
		}
	} else {
		currentPoint = ud.Point().Int()
		lim = ud.Count()
	}

	newPoint := currentPoint + addPoint

	p, err := domain.NewPoint(newPoint)
	if err != nil {
		return 0, errors.NewError("ポイントを作成できません", err)
	}

	newUserData, err := user_data.NewUserData(sID, uID, p, lim)
	if err != nil {
		return 0, errors.NewError("ユーザーデータを作成できません", err)
	}

	if err = gw.Upsert(newUserData); err != nil {
		return 0, errors.NewError("ユーザーデータの更新に失敗しました", err)
	}

	return newPoint, nil
}

// ガチャを引いたときに、カウントを+1します
//
// Todayの場合は、numをcountValue分追加
// 昨日以前の場合は、日付をTodayに変更し、numに1をセット
func IncrementCount(tx *gorm.DB, serverID, userID string, countValue int) error {
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
	currentCount := count.Count{}

	ud, err := gw.FindByIDForUpdate(id)
	if err != nil {
		if !errors.IsNotFoundError(err) {
			return errors.NewError("IDでユーザーデータを取得できません", err)
		}
	} else {
		currentPoint = ud.Point().Int()
		currentCount = ud.Count()
	}

	p, err := domain.NewPoint(currentPoint)
	if err != nil {
		return errors.NewError("ポイントを作成できません", err)
	}

	newCount := count.Count{}
	// 日付が今日の日付と一致する場合は、numをcountValue分追加
	// 今日ではない場合(昨日以前)の場合は、日付を今日にして、numに1をセット
	if currentCount.IsToday() {
		newNum, err := count.NewNum(currentCount.Num().Int() + countValue)
		if err != nil {
			return errors.NewError("回数を作成できません", err)
		}
		newCount, err = count.NewCount(currentCount.Date(), newNum)
		if err != nil {
			return errors.NewError("1日の回数を作成できません", err)
		}
	} else {
		newNum, err := count.NewNum(1)
		if err != nil {
			return errors.NewError("回数を作成できません", err)
		}
		newCount, err = count.NewCount(now.NowJST(), newNum)
		if err != nil {
			return errors.NewError("1日の回数を作成できません", err)
		}
	}

	newUserData, err := user_data.NewUserData(sID, uID, p, newCount)
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

// そのサーバーで、ポイントが上位100位のユーザーデータを取得します
func FindTop100ByServerID(tx *gorm.DB, serverID string) ([]user_data.UserData, error) {
	sID, err := domain.NewDiscordID(serverID)
	if err != nil {
		return nil, errors.NewError("サーバーIDの生成に失敗しました", err)
	}

	gw, err := gatewayUserData.NewGateway(tx)
	if err != nil {
		return nil, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	ss, err := gw.FindTop100ByServerID(sID)
	if err != nil {
		return nil, errors.NewError("サーバーIDでユーザーデータを取得できません", err)
	}

	return ss, nil
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
