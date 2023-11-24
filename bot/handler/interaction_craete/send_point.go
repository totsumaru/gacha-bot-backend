package interaction_craete

import (
	"fmt"
	"net/url"

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

	imageURL := fmt.Sprintf(
		"https://gacha-bot-image-generate.vercel.app/api/card?username=%s&avatar=%s&point=%d",
		url.QueryEscape(i.Member.User.Username),
		url.QueryEscape(i.Member.User.AvatarURL("")),
		currentPoint,
	)

	// MEMO: // URLではなくて、画像データで送りたい時はここを使用。
	// 現状はEmbedに納めたいのでコメントアウトしている。
	//
	//// 画像URLからデータを取得
	//httpRes, err := http.Get(imageURL)
	//if err != nil {
	//	return errors.NewError("画像を取得できません", err)
	//}
	//defer httpRes.Body.Close()
	//
	//imageData, err := io.ReadAll(httpRes.Body)
	//if err != nil {
	//	return errors.NewError("画像データの読み取り中にエラーが発生しました", err)
	//}
	//
	//imageReader := bytes.NewReader(imageData)

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf(description, currentPoint),
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
	}

	webhook := &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
		//Files: []*discordgo.File{{
		//	Name:   "image.jpg",
		//	Reader: imageReader,
		//}},
	}
	if _, err = editFunc(i.Interaction, webhook); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	return nil
}
