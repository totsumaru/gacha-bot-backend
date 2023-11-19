package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed/button"
)

var ButtonStyleToDiscordStyle = map[string]discordgo.ButtonStyle{
	button.ButtonStylePrimary:   discordgo.PrimaryButton,
	button.ButtonStyleSecondary: discordgo.SecondaryButton,
	button.ButtonStyleSuccess:   discordgo.SuccessButton,
	button.ButtonStyleDanger:    discordgo.DangerButton,
	button.ButtonStyleLink:      discordgo.LinkButton,
}
