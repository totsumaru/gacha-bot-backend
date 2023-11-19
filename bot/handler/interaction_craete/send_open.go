package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed/button"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

var ButtonStyleToDiscordStyle = map[string]discordgo.ButtonStyle{
	button.ButtonStylePrimary:   discordgo.PrimaryButton,
	button.ButtonStyleSecondary: discordgo.SecondaryButton,
	button.ButtonStyleSuccess:   discordgo.SuccessButton,
	button.ButtonStyleDanger:    discordgo.DangerButton,
	button.ButtonStyleLink:      discordgo.LinkButton,
}

// Openのメッセージを送信します
func SendOpen(s *discordgo.Session, i *discordgo.InteractionCreate, domainGacha gacha.Gacha) error {
	editFunc, err := SendInteractionWaitingMessage(s, i, false, true)
	if err != nil {
		return errors.NewError("Waitingメッセージが送信できません")
	}

	btn1 := discordgo.Button{
		Label:    domainGacha.Open().Button()[0].Label().String(),
		Style:    ButtonStyleToDiscordStyle[domainGacha.Open().Button()[0].Style().String()],
		CustomID: domainGacha.Open().Button()[0].Kind().String(),
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn1},
	}

	embed := &discordgo.MessageEmbed{
		Description: domainGacha.Open().Description().String(),
		Image: &discordgo.MessageEmbedImage{
			URL: domainGacha.Open().ImageURL().String(),
		},
		Color: domainGacha.Open().Color().Int(),
	}

	webhook := &discordgo.WebhookEdit{
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &[]discordgo.MessageComponent{actions},
	}
	if _, err = editFunc(i.Interaction, webhook); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}
