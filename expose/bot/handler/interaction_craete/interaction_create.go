package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/expose"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// インタラクションが作成された時のハンドラです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := expose.DB.Transaction(func(tx *gorm.DB) error {
		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), i.GuildID)
		return
	}
}
