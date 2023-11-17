package upsert

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiGacha "github.com/totsumaru/gacha-bot-backend/api/gacha"
	"github.com/totsumaru/gacha-bot-backend/app/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// ガチャを作成/更新します
func UpsertGacha(e *gin.Engine, db *gorm.DB) {
	e.POST("/gacha/upsert", func(c *gin.Context) {
		gachaID := c.Query("gacha_id")
		//authHeader := c.GetHeader(auth.HeaderAuthorization)
		//
		//var userID string

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
		//	if err = auth.IsAdmin(serverID, userID); err != nil {
		//		errors.HandleError(c, 401, "管理者ではありません", err)
		//		return
		//	}
		//}

		var gachaReq apiGacha.GachaReq
		// リクエストボディのJSONをGachaReqにバインド
		if err := c.ShouldBindJSON(&gachaReq); err != nil {
			errors.HandleError(c, http.StatusBadRequest, "リクエストの解析に失敗しました", err)
			return
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			appReq := apiGacha.ConvertToAppGachaReq(gachaReq)
			_, err := gacha.UpsertGacha(tx, gachaID, appReq)
			if err != nil {
				return errors.NewError("ガチャの更新に失敗しました", err)
			}

			return nil
		})
		if err != nil {
			errors.HandleError(c, 500, "ガチャを更新できません", err)
			return
		}

		c.JSON(http.StatusOK, "")
	})
}
