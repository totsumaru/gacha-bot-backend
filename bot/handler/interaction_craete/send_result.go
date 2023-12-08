package interaction_craete

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	appUserData "github.com/totsumaru/gacha-bot-backend/application/user_data"
	domainGacha "github.com/totsumaru/gacha-bot-backend/domain/gacha"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/result"
	"github.com/totsumaru/gacha-bot-backend/lib/color"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// 結果メッセージを送信します
//
// 追加された後の最新のポイントを返します。
func SendResult(
	tx *gorm.DB,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	domainGacha domainGacha.Gacha,
) (int, error) {
	editFunc, err := SendInteractionWaitingMessage(s, i, true, true)
	if err != nil {
		return 0, errors.NewError("Waitingメッセージが送信できません")
	}

	// 現在のカウントを取得する
	ud, err := appUserData.FindByServerIDAndUserID(tx, i.GuildID, i.Member.User.ID)
	if err != nil && !errors.IsNotFoundError(err) {
		return 0, errors.NewError("ユーザーデータを取得できません", err)
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
			return 0, errors.NewError("レスポンスを送信できません", err)
		}
		return 0, nil
	}

	r := chooseProb(domainGacha.Result())

	embed := &discordgo.MessageEmbed{
		Description: r.Embed().Description().String(),
		Image: &discordgo.MessageEmbedImage{
			URL: r.Embed().ImageURL().String(),
		},
		Color: r.Embed().Color().Int(),
	}

	webhook := &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}

	// ボタンを追加します
	// MEMO: Linkに限定しています
	components := make([]discordgo.MessageComponent, 0)
	for _, btn := range r.Embed().Button() {
		if btn.IsVisible() {
			components = append(components, discordgo.Button{
				Label: btn.Label().String(),
				Style: discordgo.LinkButton,
				URL:   btn.URL().String(),
			})
		}
	}
	if len(components) > 0 {
		actions := discordgo.ActionsRow{
			Components: components,
		}
		webhook.Components = &[]discordgo.MessageComponent{actions}
	}

	if _, err = editFunc(i.Interaction, webhook); err != nil {
		return 0, errors.NewError("レスポンスを送信できません", err)
	}

	// ポイントを追加
	latestPoint, err := appUserData.AddPoint(tx, i.GuildID, i.Member.User.ID, r.Point().Int())
	if err != nil {
		return 0, errors.NewError("ポイントを追加できません", err)
	}
	// カウントを追加
	if err = appUserData.IncrementCount(tx, i.GuildID, i.Member.User.ID, 1); err != nil {
		return 0, errors.NewError("カウントを追加できません", err)
	}

	return latestPoint, nil
}

// resultを確率に従って1つ取得します
func chooseProb(results []result.Result) result.Result {
	// 新しい乱数生成器を作成
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// 累積確率リストを作成
	cumulativeProbs := make([]int, len(results))
	total := 0
	for i, rslt := range results {
		total += rslt.Probability().Int()
		cumulativeProbs[i] = total
	}

	// 0 から合計確率までの乱数を生成
	randVal := r.Intn(total)

	// 乱数が累積確率の範囲内に入る最初の要素を選択
	for i, cumulativeProb := range cumulativeProbs {
		if randVal < cumulativeProb {
			return results[i]
		}
	}

	// 通常はここに到達しないが、念のため最後の要素を返す
	return results[len(results)-1]
}
