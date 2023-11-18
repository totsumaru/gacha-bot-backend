package gacha

import (
	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	gatewayGacha "github.com/totsumaru/gacha-bot-backend/gateway/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// サーバーIDでガチャを取得します
//
// 現状は1つなので問題ありません。
func FindByServerID(tx *gorm.DB, serverID string) (gacha.Gacha, error) {
	svID, err := domain.NewDiscordID(serverID)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("サーバーIDを作成できません", err)
	}

	gw, err := gatewayGacha.NewGateway(tx)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	s, err := gw.FindByServerID(svID)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("サーバーIDでガチャを取得できません", err)
	}

	return s, nil
}
