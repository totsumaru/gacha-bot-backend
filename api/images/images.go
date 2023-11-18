package images

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/totsumaru/gacha-bot-backend/application/cloudflare"
	"github.com/totsumaru/gacha-bot-backend/lib/auth"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// 画像をアップロードします
func UploadImages(e *gin.Engine) {
	e.POST("/api/images", func(c *gin.Context) {
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

		form, err := c.MultipartForm()
		if err != nil {
			log.Println("フォームデータの取得に失敗しました", err)
			c.JSON(500, gin.H{"error": "フォームデータの取得に失敗しました"})
			return
		}

		images := form.File["images"] // "images" はフロントエンドで設定されたフィールド名

		var urls []string
		for _, image := range images {
			url, err := cloudflare.Upload(image)
			if err != nil {
				log.Println("画像のアップロードに失敗しました", err)
				c.JSON(500, gin.H{"error": "画像のアップロードに失敗しました"})
				return
			}
			urls = append(urls, url)
		}

		c.JSON(200, gin.H{"urls": urls})
	})
}
