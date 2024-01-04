package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	appUserData "github.com/totsumaru/gacha-bot-backend/application/user_data"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/color"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// Openのメッセージを送信します
func SendOpen(
	tx *gorm.DB,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	domainGacha gacha.Gacha,
) error {
	editFunc, err := SendInteractionWaitingMessage(s, i, false, true)
	if err != nil {
		return errors.NewError("Waitingメッセージが送信できません")
	}

	// 現在のカウントを取得する
	ud, err := appUserData.FindByServerIDAndUserID(tx, i.GuildID, i.Member.User.ID)
	if err != nil && !errors.IsNotFoundError(err) {
		return errors.NewError("ユーザーデータを取得できません", err)
	}
	// カウントが今日かつ、1回以上の場合は、エラーメッセージを送信する
	if ud.Count().IsToday() && ud.Count().Num().Int() > 0 {
		embed := &discordgo.MessageEmbed{
			Description: "1日の上限回数に達しています。\n明日になってからもう一度お試しください。",
			Color:       color.Red,
		}

		webhook := &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		}
		if _, err = editFunc(i.Interaction, webhook); err != nil {
			return errors.NewError("レスポンスを送信できません", err)
		}
		return nil
	}

	btn1 := discordgo.Button{
		Label:    domainGacha.Open().Button()[0].Label().String(),
		Style:    ButtonStyleToDiscordStyle[domainGacha.Open().Button()[0].Style().String()],
		CustomID: domainGacha.Open().Button()[0].Kind().String(),
		Emoji:    discordgo.ComponentEmoji{Name: "▶️"},
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
