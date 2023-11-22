package interaction_craete

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	appGacha "github.com/totsumaru/gacha-bot-backend/application/gacha"
	"github.com/totsumaru/gacha-bot-backend/application/user_data"
	domainGacha "github.com/totsumaru/gacha-bot-backend/domain/gacha"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/result"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// 結果メッセージを送信します
func SendResult(
	tx *gorm.DB,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	domainGacha domainGacha.Gacha,
) error {
	editFunc, err := SendInteractionWaitingMessage(s, i, true, true)
	if err != nil {
		return errors.NewError("Waitingメッセージが送信できません")
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
	if _, err = editFunc(i.Interaction, webhook); err != nil {
		return errors.NewError("レスポンスを送信できません", err)
	}

	// ポイントを追加
	if err = user_data.AddPoint(
		tx, i.GuildID, i.Member.User.ID, r.Point().Int(),
	); err != nil {
		return errors.NewError("ポイントを追加できません", err)
	}
	// カウントを追加
	if err = user_data.IncrementCount(tx, i.GuildID, i.Member.User.ID); err != nil {
		return errors.NewError("カウントを追加できません", err)
	}

	// ポイントがロール付与の条件を満たしていた場合は、ロールを付与
	{
		ud, err := user_data.FindByServerIDAndUserID(tx, i.GuildID, i.Member.User.ID)
		if err != nil && !errors.IsNotFoundError(err) {
			return errors.NewError("ユーザーデータを取得できません", err)
		}

		totalPoint := ud.Point().Int() + r.Point().Int()

		// ガチャ情報を取得
		ga, err := appGacha.FindByServerID(tx, i.GuildID)
		if err != nil {
			return errors.NewError("ガチャを取得できません", err)
		}

		// ロール付与の条件を満たしているか確認
		for _, ro := range ga.Role() {
			if totalPoint >= ro.Point().Int() {
				// 指定のロールを持っているかを確認します
				hasRole := false
				for _, mr := range i.Member.Roles {
					if mr == ro.ID().String() {
						hasRole = true
						break
					}
				}

				// 指定のロールを持っていない場合は、ロールを付与します
				if !hasRole {
					if err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, ro.ID().String()); err != nil {
						return errors.NewError("ロールを付与できません", err)
					}
				}
			}
		}
	}

	return nil
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
