package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/totsumaru/gacha-bot-backend/app/server"
	domainServer "github.com/totsumaru/gacha-bot-backend/domain/server"
	"github.com/totsumaru/gacha-bot-backend/lib/auth"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// レスポンスです
type Res struct {
	ID           string `json:"id"`
	SubscriberID string `json:"subscriber_id"`
	Stripe       struct {
		CustomerID     string `json:"customer_id"`
		SubscriptionID string `json:"subscription_id"`
	} `json:"stripe"`
	IsSubscriber bool `json:"is_subscriber"`
}

// サーバーの情報を取得します
func GetServer(e *gin.Engine, db *gorm.DB) {
	e.GET("/api/server", func(c *gin.Context) {
		serverID := c.Query("server_id")
		authHeader := c.GetHeader(auth.HeaderAuthorization)

		var userDiscordID string

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
			userDiscordID = headerRes.DiscordID

			if err = auth.IsAdmin(serverID, userDiscordID); err != nil {
				errors.HandleError(c, 401, "管理者ではありません", err)
				return
			}
		}

		res := Res{}
		err := db.Transaction(func(tx *gorm.DB) error {
			sv, err := server.FindByID(tx, serverID)
			if err != nil {
				return err
			}

			res = convertDomainServerToApiRes(sv, userDiscordID)

			return nil
		})
		if err != nil {
			errors.HandleError(c, 500, "サーバーの情報を取得できません", err)
			return
		}

		c.JSON(http.StatusOK, res)
	})
}

// サーバーのドメインをAPIのレスポンスに変換します
func convertDomainServerToApiRes(domainSv domainServer.Server, userDiscordID string) Res {
	res := Res{}
	res.ID = domainSv.ID().String()
	res.SubscriberID = domainSv.SubscriberID().String()
	res.Stripe.CustomerID = domainSv.Stripe().CustomerID().String()
	res.Stripe.SubscriptionID = domainSv.Stripe().SubscriptionID().String()
	res.IsSubscriber = domainSv.SubscriberID().String() == userDiscordID

	return res
}
