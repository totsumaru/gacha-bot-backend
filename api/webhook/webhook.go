package webhook

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/totsumaru/gacha-bot-backend/application/server"
	"github.com/totsumaru/gacha-bot-backend/lib/discord"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// Stripeからコールされるwebhookです
func Webhook(e *gin.Engine, db *gorm.DB) {
	e.POST("/api/webhook", func(c *gin.Context) {
		header := c.GetHeader("Stripe-Signature")
		webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			errors.HandleError(c, 500, "bodyを読み取れません", err)
			return
		}

		// verify
		event, err := webhook.ConstructEvent(body, header, webhookSecret)
		if err != nil {
			errors.HandleError(c, 400, "リクエストが不正です", err)
			return
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			// イベントオブジェクトのdocument
			switch event.Type {
			case "checkout.session.completed":
				// Checkout で顧客が「支払う」または「登録」ボタンをクリックすると送信され、新しい購入が通知されます。
				// `CUSTOMER_ID`をDBに登録する必要があります。
				customerID := event.Data.Object["customer"].(string)
				subscriptionID := event.Data.Object["subscription"].(string)
				metadata := event.Data.Object["metadata"].(map[string]interface{})
				guildID := metadata["guild_id"].(string)
				discordID := metadata["discord_id"].(string)

				if err = server.StartSubscription(
					tx, guildID, discordID, customerID, subscriptionID,
				); err != nil {
					return errors.NewError("サブスクリプションを開始できません", err)
				}
			case "invoice.paid":
				// MEMO: 何を実装すべきか不明
				// 請求期間ごとに、支払いが成功すると送信されるイベントです。
				// 初回の支払い時にもコールされます。
			case "customer.subscription.deleted":
				// 顧客のサブスクリプションが終了すると送信されます。
				metadata := event.Data.Object["metadata"].(map[string]string)
				guildID := metadata["guild_id"]

				if err = server.DeleteSubscription(tx, guildID); err != nil {
					return errors.NewError("サブスクリプションを削除できません", err)
				}
			case "invoice.payment_failed":
				customerID := event.Data.Object["customer"].(string)
				subscriptionID := event.Data.Object["id"].(string)
				metadata := event.Data.Object["metadata"].(map[string]string)
				guildID := metadata["guild_id"]
				discordID := metadata["discord_id"]

				resObj := map[string]string{
					"customerID":     customerID,
					"subscriptionID": subscriptionID,
					"guildID":        guildID,
					"discordID":      discordID,
				}

				errors.SendErrMsg(
					discord.Session,
					fmt.Errorf("サブスクリプションの支払いに失敗したユーザーがいます。: %v", resObj),
					guildID,
				)

				return errors.NewError(fmt.Sprintf(
					"サブスクリプションの支払いに失敗したユーザーがいます。: %v",
					resObj,
				))
			default:
			}
			return nil
		})
		if err != nil {
			log.Println("エラーが発生しました: ", err)
			return
		}

		c.JSON(http.StatusOK, "")
	})
}
