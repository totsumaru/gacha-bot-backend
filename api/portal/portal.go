package portal

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/billingportal/session"
	"github.com/totsumaru/gacha-bot-backend/app/server"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// レスポンスです
type Res struct {
	RedirectURL string `json:"redirect_url"`
}

// カスタマーポータルのURLを作成します
func CreateCustomerPortal(e *gin.Engine, db *gorm.DB) {
	// カスタマーポータルのURLを作成します
	e.POST("/api/portal", func(c *gin.Context) {
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		serverID := c.Query("server_id")
		//authHeader := c.GetHeader(auth.HeaderAuthorization)

		var userID string
		userID = "960104306151948328" // TODO: テスト用のため削除する

		// verify
		//{
		//	if serverID == "" || authHeader == "" {
		//		errors.HandleError(c, 400, "リクエストが不正です", fmt.Errorf(
		//			"serverID: %s, authHeader: %s", serverID, authHeader,
		//		))
		//		return
		//	}
		//
		//	headerRes, err := auth.GetAuthHeader(authHeader)
		//	if err != nil {
		//		errors.HandleError(c, 401, "トークンの認証に失敗しました", err)
		//		return
		//	}
		//	userID = headerRes.DiscordID
		//
		//	if err = auth.IsAdmin(serverID, headerRes.DiscordID); err != nil {
		//		errors.HandleError(c, 401, "必要な権限を持っていません", err)
		//		return
		//	}
		//}

		var url string
		err := db.Transaction(func(tx *gorm.DB) error {
			// そのサーバーの支払い者が本人であるかを確認
			apiRes, err := server.FindByID(tx, serverID)
			if err != nil {
				return errors.NewError("IDでサーバーを取得できません", err)
			}
			if apiRes.SubscriberID().String() != userID {
				return errors.NewError("本人ではありません", err)
			}

			// カスタマーポータルから戻るボタンで、遷移するURLを作成
			returnURL := fmt.Sprintf(
				"%s/dashboard/%s/config",
				os.Getenv("FRONTEND_URL"),
				serverID,
			)

			// customerIdからカスタマーポータルURLを作成
			params := &stripe.BillingPortalSessionParams{
				Customer:  stripe.String(apiRes.Stripe().CustomerID().Value()),
				ReturnURL: stripe.String(returnURL),
			}

			s, err := session.New(params)
			if err != nil {
				return errors.NewError("stripeのsessionを作成できません", err)
			}
			url = s.URL

			return nil
		})
		if err != nil {
			errors.HandleError(c, 500, "stripeのsessionを作成できません", err)
			return
		}

		c.JSON(http.StatusOK, Res{
			RedirectURL: url,
		})
	})
}
