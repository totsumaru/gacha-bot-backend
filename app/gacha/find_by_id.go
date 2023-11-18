package gacha

import (
	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	gatewayGacha "github.com/totsumaru/gacha-bot-backend/gateway/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// IDでガチャを取得します
func FindByID(tx *gorm.DB, id string) (gacha.Gacha, error) {
	gachaID, err := domain.RestoreUUID(id)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("IDを作成できません", err)
	}

	gw, err := gatewayGacha.NewGateway(tx)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	s, err := gw.FindByID(gachaID)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("IDでサーバーを取得できません", err)
	}

	return s, nil
}
