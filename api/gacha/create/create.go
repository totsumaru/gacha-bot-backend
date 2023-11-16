package create

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiGacha "github.com/totsumaru/gacha-bot-backend/api/gacha"
	"github.com/totsumaru/gacha-bot-backend/app/gacha"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

type Req struct {
	apiGacha.GachaReq
}

// ガチャを作成します
func CreateGacha(e *gin.Engine, db *gorm.DB) {
	e.POST("/gacha/create", func(c *gin.Context) {
		//authHeader := c.GetHeader(auth.HeaderAuthorization)
		//
		//var userID string

		//// verify
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

		gachaReq := apiGacha.GachaReq{}
		// リクエストボディのJSONをGachaReqにバインド
		if err := c.ShouldBindJSON(&gachaReq); err != nil {
			errors.HandleError(c, http.StatusBadRequest, "リクエストの解析に失敗しました", err)
			return
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			appReq := apiGacha.ConvertToAppGachaReq(gachaReq)
			_, err := gacha.CreateGacha(tx, appReq)
			if err != nil {
				return errors.NewError("ガチャの作成に失敗しました", err)
			}

			return nil
		})
		if err != nil {
			errors.HandleError(c, 500, "ガチャを作成できません", err)
			return
		}

		c.JSON(http.StatusOK, "")
	})
}
