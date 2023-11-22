package interaction_craete

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/application/gacha"
	"github.com/totsumaru/gacha-bot-backend/application/user_data"
	"github.com/totsumaru/gacha-bot-backend/bot"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed/button"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// インタラクションが作成された時のハンドラです
func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		// パネルのボタンが押された時、Openのメッセージを送信します
		case button.ButtonKindToOpen:
			err := bot.DB.Transaction(func(tx *gorm.DB) error {
				ga, err := gacha.FindByServerID(tx, i.GuildID)
				if err != nil {
					return errors.NewError("ガチャを取得できません", err)
				}

				if err = SendOpen(tx, s, i, ga); err != nil {
					return errors.NewError("Openのメッセージを送信できません", err)
				}
				return nil
			})
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), i.GuildID)
				return
			}
		case button.ButtonKindToResult:
			totalPoint := 0
			err := bot.DB.Transaction(func(tx *gorm.DB) error {
				// TODO: ForUpdateを使う
				ga, err := gacha.FindByServerID(tx, i.GuildID)
				if err != nil {
					return errors.NewError("ガチャを取得できません", err)
				}

				totalPoint, err = SendResult(tx, s, i, ga)
				if err != nil {
					return errors.NewError("結果を送信できません", err)
				}

				return nil
			})
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), i.GuildID)
				return
			}

			// ポイントがロール付与の条件を満たしていた場合は、ロールを付与
			err = bot.DB.Transaction(func(tx *gorm.DB) error {
				// ガチャ情報を取得
				ga, err := gacha.FindByServerID(tx, i.GuildID)
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

				return nil
			})
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), i.GuildID)
				return
			}
		case "check_point":
			err := bot.DB.Transaction(func(tx *gorm.DB) error {
				ud, err := user_data.FindByServerIDAndUserID(tx, i.GuildID, i.Member.User.ID)
				if err != nil && !errors.IsNotFoundError(err) {
					return errors.NewError("ユーザーデータを取得できません", err)
				}

				if err = SendPoint(s, i, ud.Point().Int()); err != nil {
					return errors.NewError("ポイントを送信できません", err)
				}
				return nil
			})
			if err != nil {
				errors.SendErrMsg(s, errors.NewError("エラーが発生しました", err), i.GuildID)
				return
			}
		}
	}
}
