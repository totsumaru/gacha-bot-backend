package gacha

import (
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	gatewayGacha "github.com/totsumaru/gacha-bot-backend/gateway/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// ガチャを新規作成します
func CreateGacha(tx *gorm.DB, req GachaReq) (gacha.Gacha, error) {
	panel, err := createEmbed(req.Panel)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("panelの生成に失敗しました", err)
	}

	open, err := createEmbed(req.Open)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("openの生成に失敗しました", err)
	}

	result, err := createResult(req.Result)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("resultの生成に失敗しました", err)
	}

	g, err := gacha.NewGacha(panel, open, result)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("ガチャの生成に失敗しました", err)
	}

	gw, err := gatewayGacha.NewGateway(tx)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("Gatewayの生成に失敗しました", err)
	}

	if err = gw.Create(g); err != nil {
		return gacha.Gacha{}, errors.NewError("ガチャの保存に失敗しました", err)
	}

	return g, nil
}
