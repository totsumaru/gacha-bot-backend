package api

import (
	"github.com/gin-gonic/gin"
	"github.com/totsumaru/gacha-bot-backend/api/admin/servers"
	"github.com/totsumaru/gacha-bot-backend/api/checkout"
	"github.com/totsumaru/gacha-bot-backend/api/gacha"
	"github.com/totsumaru/gacha-bot-backend/api/gacha/ranking"
	"github.com/totsumaru/gacha-bot-backend/api/gacha/upsert"
	"github.com/totsumaru/gacha-bot-backend/api/images"
	"github.com/totsumaru/gacha-bot-backend/api/portal"
	"github.com/totsumaru/gacha-bot-backend/api/server"
	"github.com/totsumaru/gacha-bot-backend/api/webhook"
	"gorm.io/gorm"
)

// ルートを設定します
func RegisterRouter(e *gin.Engine, db *gorm.DB) {
	Route(e)
	checkout.Checkout(e)
	portal.CreateCustomerPortal(e, db)
	webhook.Webhook(e, db)
	upsert.UpsertGacha(e, db)
	images.UploadImages(e)
	gacha.GetGacha(e, db)
	server.GetServer(e, db)
	ranking.GetRanking(e, db)
	servers.GetAdminServers(e, db)
}

// ルートです
//
// Note: この関数は削除しても問題ありません
func Route(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
}
