package gacha

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/totsumaru/gacha-bot-backend/app/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/auth"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// ガチャを取得します
func GetGacha(e *gin.Engine, db *gorm.DB) {
	e.GET("/api/gacha", func(c *gin.Context) {
		serverID := c.Query("server_id")
		authHeader := c.GetHeader(auth.HeaderAuthorization)

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

			if err = auth.IsAdmin(serverID, userID); err != nil {
				errors.HandleError(c, 401, "管理者ではありません", err)
				return
			}
		}

		res := GachaRes{}
		err := db.Transaction(func(tx *gorm.DB) error {
			ga, err := gacha.FindByServerID(tx, serverID)
			if err != nil {
				return err
			}

			res = ConvertToAPIGachaRes(ga)

			return nil
		})
		if err != nil {
			if errors.IsNotFoundError(err) {
				res.ID = uuid.NewString()
				res.ServerID = serverID
				res.Panel.Button = []ButtonReq{createPanelBtn()}
				res.Open.Button = []ButtonReq{createOpenBtn()}
				res.Result = []ResultReq{{
					Embed:       EmbedReq{},
					Point:       0,
					Probability: 100,
				}}

				c.JSON(http.StatusOK, res)
				return
			}
			errors.HandleError(c, 500, "ガチャを取得できません", err)
			return
		}

		c.JSON(http.StatusOK, res)
	})
}

// panelのボタンを作成します
func createPanelBtn() ButtonReq {
	return ButtonReq{
		Kind:  "to_open",
		Label: "ガチャを回す",
		Style: "PRIMARY",
	}
}

// openのボタンを作成します
func createOpenBtn() ButtonReq {
	return ButtonReq{
		Kind:  "to_result",
		Label: "結果を見る",
		Style: "PRIMARY",
	}
}
