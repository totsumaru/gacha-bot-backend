package gacha

import (
	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	gatewayGacha "github.com/totsumaru/gacha-bot-backend/gateway/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// ガチャをUpsertします
func UpsertGacha(tx *gorm.DB, req GachaReq) (gacha.Gacha, error) {
	i, err := domain.RestoreUUID(req.ID)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("DiscordIDの生成に失敗しました", err)
	}

	sID, err := domain.NewDiscordID(req.ServerID)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("サーバーIDの生成に失敗しました", err)
	}

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

	g, err := gacha.RestoreGacha(i, sID, panel, open, result)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("ガチャの生成に失敗しました", err)
	}

	gw, err := gatewayGacha.NewGateway(tx)
	if err != nil {
		return gacha.Gacha{}, errors.NewError("Gatewayの生成に失敗しました", err)
	}

	if err = gw.Upsert(g); err != nil {
		return gacha.Gacha{}, errors.NewError("ガチャの更新に失敗しました", err)
	}

	return g, nil
}
