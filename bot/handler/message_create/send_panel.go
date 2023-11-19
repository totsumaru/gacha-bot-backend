package message_create

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// パネルメッセージを送信します
func SendPanel(s *discordgo.Session, m *discordgo.MessageCreate, gachaDomain gacha.Gacha) error {
	btn1 := discordgo.Button{
		Label:    "ガチャを回す！",
		Style:    discordgo.PrimaryButton,
		CustomID: gachaDomain.Panel().Button()[0].Kind().String(),
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
	}

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: gachaDomain.Panel().ImageURL().String(),
		},
		Title:       "ロールガチャ",
		Description: gachaDomain.Panel().Description().String(),
		Color:       gachaDomain.Panel().Color().Int(),
	}

	// 新規のパネルを作成します
	messageSend := &discordgo.MessageSend{
		Components: []discordgo.MessageComponent{actions},
		Embed:      embed,
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, messageSend)
	if err != nil {
		return errors.NewError("パネルメッセージを送信できません", err)
	}

	return nil
}