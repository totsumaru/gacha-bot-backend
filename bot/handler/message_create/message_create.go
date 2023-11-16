package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/app/server"
	"github.com/totsumaru/gacha-bot-backend/bot"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// メッセージが作成された時のハンドラです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := bot.DB.Transaction(func(tx *gorm.DB) error {
		switch m.Content {
		case "!gacha-setup":
			// すでに登録されている場合は、返信を返すのみ
			_, err := server.FindByID(tx, m.GuildID)
			if err == nil {
				if _, err = s.ChannelMessageSend(m.ChannelID, "すでに登録されています"); err != nil {
					return errors.NewError("メッセージの送信に失敗しました", err)
				}
				return nil
			}

			if err = server.CreateServer(tx, m.GuildID); err != nil {
				return errors.NewError("サーバーの作成に失敗しました", err)
			}

			if _, err = s.ChannelMessageSend(m.ChannelID, "セットアップが完了しました"); err != nil {
				return errors.NewError("メッセージの送信に失敗しました", err)
			}
		}
		return nil
	})
	if err != nil {
		errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), m.GuildID)
		return
	}
}
