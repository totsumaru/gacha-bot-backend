package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/expose"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// メッセージが作成された時のハンドラです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := expose.DB.Transaction(func(tx *gorm.DB) error {
		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), m.GuildID)
		return
	}
}
