package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// レスポンスのEdit関数を返します
type EditFunc func(interaction *discordgo.Interaction, newresp *discordgo.WebhookEdit, options ...discordgo.RequestOption) (*discordgo.Message, error)

// Interactionのdeferメッセージを送信します
func SendInteractionWaitingMessage(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	isUpdate bool,
	isEphemeral bool,
) (EditFunc, error) {
	responseType := discordgo.InteractionResponseDeferredChannelMessageWithSource
	if isUpdate {
		responseType = discordgo.InteractionResponseUpdateMessage
	}

	embed := &discordgo.MessageEmbed{
		Description: "Loading...",
	}

	resp := &discordgo.InteractionResponse{
		Type: responseType,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		}, // isEphemeralの値を入れるためDataのフィールドは定義しておく
	}

	if isEphemeral {
		resp.Data.Flags = discordgo.MessageFlagsEphemeral
	}

	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		return nil, errors.NewError("インタラクションのレスポンスを送信できません", err)
	}

	return s.InteractionResponseEdit, nil
}
