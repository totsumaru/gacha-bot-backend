package message_create

import (
	"github.com/bwmarrin/discordgo"
	appUserData "github.com/totsumaru/gacha-bot-backend/application/user_data"
	"github.com/totsumaru/gacha-bot-backend/bot"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// 回数をリセットします
func ResetNum(s *discordgo.Session, m *discordgo.MessageCreate) error {
	ud, err := appUserData.FindByServerIDAndUserIDForUpdate(bot.DB, m.GuildID, m.Author.ID)
	if err != nil {
		return errors.NewError("ユーザーデータの取得に失敗しました", err)
	}

	// 現在のカウント分をマイナスします
	if err = appUserData.IncrementCount(
		bot.DB,
		m.GuildID,
		m.Author.ID,
		-ud.Count().Num().Int(),
	); err != nil {
		return errors.NewError("カウントを追加できません", err)
	}

	// 結果を送信します
	_, err = s.ChannelMessageSend(m.ChannelID, "回数をリセットしました")
	if err != nil {
		return errors.NewError("メッセージの送信に失敗しました", err)
	}

	return nil
}
