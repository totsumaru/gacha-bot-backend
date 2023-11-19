package message_create

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/application/gacha"
	"github.com/totsumaru/gacha-bot-backend/application/server"
	"github.com/totsumaru/gacha-bot-backend/bot"
	"github.com/totsumaru/gacha-bot-backend/lib/auth"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// メッセージが作成された時のハンドラです
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case "!gacha-setup":
		// 管理者以外は実行できません
		if auth.IsAdmin(m.GuildID, m.Author.ID) != nil {
			if _, err := s.ChannelMessageSend(
				m.ChannelID,
				"このコマンドは管理者のみ実行できます",
			); err != nil {
				errors.SendErrMsg(s, errors.NewError("メッセージの送信に失敗しました", err), m.GuildID)
				return
			}
			return
		}
		err := bot.DB.Transaction(func(tx *gorm.DB) error {
			dashboardURL := os.Getenv("FRONTEND_URL") + "/server/" + m.GuildID
			// すでに登録されている場合は、返信を返すのみ
			_, err := server.FindByID(tx, m.GuildID)
			if err == nil {
				if _, err = s.ChannelMessageSend(
					m.ChannelID,
					fmt.Sprintf("すでに登録されています。以下のURLから設定を進めてください。\n%s", dashboardURL),
				); err != nil {
					return errors.NewError("メッセージの送信に失敗しました", err)
				}
				return nil
			}

			if err = server.CreateServer(tx, m.GuildID); err != nil {
				return errors.NewError("サーバーの作成に失敗しました", err)
			}

			if _, err = s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("セットアップが完了しました。以下のURLから設定を進めてください。\n%s", dashboardURL),
			); err != nil {
				return errors.NewError("メッセージの送信に失敗しました", err)
			}
			return nil
		})
		if err != nil {
			errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), m.GuildID)
			return
		}
	case "!gacha-panel":
		if auth.IsAdmin(m.GuildID, m.Author.ID) != nil {
			if _, err := s.ChannelMessageSend(
				m.ChannelID,
				"このコマンドは管理者のみ実行できます",
			); err != nil {
				errors.SendErrMsg(s, errors.NewError("メッセージの送信に失敗しました", err), m.GuildID)
				return
			}
			return
		}
		err := bot.DB.Transaction(func(tx *gorm.DB) error {
			ga, err := gacha.FindByServerID(tx, m.GuildID)
			if err != nil {
				if errors.IsNotFoundError(err) {
					if _, err = s.ChannelMessageSend(
						m.ChannelID,
						"情報が登録されていません。「!gacha-setup」を実行して登録してください",
					); err != nil {
						return errors.NewError("メッセージの送信に失敗しました", err)
					}
					return nil
				}
				return errors.NewError("ガチャの取得に失敗しました", err)
			}

			if err = SendPanel(s, m, ga); err != nil {
				return errors.NewError("パネルメッセージの送信に失敗しました", err)
			}
			return nil
		})
		if err != nil {
			errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), m.GuildID)
			return
		}
	}
}
