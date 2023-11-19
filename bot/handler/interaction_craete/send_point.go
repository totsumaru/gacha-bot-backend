package interaction_craete

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ポイントのメッセージを送信します
func SendPoint(s *discordgo.Session, i *discordgo.InteractionCreate, currentPoint int) error {
	editFunc, err := SendInteractionWaitingMessage(s, i, false, true)
	if err != nil {
		return errors.NewError("Waitingメッセージが送信できません")
	}

	description := `
現在のポイント

**%d** pt
`

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(description, currentPoint),
	}

	webhook := &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}
	if _, err = editFunc(i.Interaction, webhook); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}
