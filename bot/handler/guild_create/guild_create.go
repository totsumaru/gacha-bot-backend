package guild_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/app/server"
	"github.com/totsumaru/gacha-bot-backend/bot"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// ギルドが作成された時のハンドラです
func GuildCreateHandler(s *discordgo.Session, i *discordgo.GuildCreate) {
	err := bot.DB.Transaction(func(tx *gorm.DB) error {
		if err := server.CreateServer(tx, i.Guild.ID); err != nil {
			return errors.NewError("サーバーの作成に失敗しました", err)
		}
		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), i.Guild.ID)
		return
	}
}
