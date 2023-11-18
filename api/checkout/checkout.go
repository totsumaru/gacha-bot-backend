package checkout

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/totsumaru/gacha-bot-backend/lib/auth"
	"github.com/totsumaru/gacha-bot-backend/lib/discord"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// レスポンスです
type Res struct {
	RedirectURL string `json:"redirect_url"`
}

func Checkout(e *gin.Engine) {
	// チェックアウトを作成します
	e.POST("/api/checkout", func(c *gin.Context) {
		authHeader := c.GetHeader(auth.HeaderAuthorization)
		serverID := c.Query("server_id")

		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		priceId := os.Getenv("STRIPE_PRICE_ID")

		var userID string

		// verify
		{
			if serverID == "" || authHeader == "" {
				errors.HandleError(c, 400, "リクエストが不正です", fmt.Errorf(
					"serverID: %s, authHeader: %s", serverID, authHeader,
				))
				return
			}

			headerRes, err := auth.GetAuthHeader(authHeader)
			if err != nil {
				errors.HandleError(c, 401, "トークンの認証に失敗しました", err)
				return
			}
			userID = headerRes.DiscordID

			if err = auth.IsAdmin(serverID, headerRes.DiscordID); err != nil {
				errors.HandleError(c, 401, "必要な権限を持っていません", err)
				return
			}
		}

		ds := discord.Session
		guild, err := ds.Guild(serverID)
		if err != nil {
			errors.HandleError(c, 500, "サーバー情報を取得できません", err)
			return
		}

		u, err := ds.User(userID)
		if err != nil {
			errors.HandleError(c, 500, "ユーザー情報を取得できません", err)
			return
		}

		FEURL := os.Getenv("FRONTEND_URL")

		successURL := fmt.Sprintf("%s/dashboard/%s/success", FEURL, serverID)
		cancelURL := fmt.Sprintf("%s/dashboard/%s/config", FEURL, serverID)

		params := &stripe.CheckoutSessionParams{
			SuccessURL: stripe.String(successURL),
			CancelURL:  stripe.String(cancelURL),
			Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(priceId),
					Quantity: stripe.Int64(1),
				},
			},
		}
		params.AddMetadata("guild_id", guild.ID)
		params.AddMetadata("discord_id", u.ID)

		s, err := session.New(params)
		if err != nil {
			errors.HandleError(c, 500, "stripeのセッションを作成できません", err)
			return
		}

		c.JSON(http.StatusOK, Res{
			RedirectURL: s.URL,
		})
	})
}
