package images

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/totsumaru/gacha-bot-backend/app/cloudflare"
)

// 画像をアップロードします
func UploadImages(e *gin.Engine) {
	e.POST("/api/images", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			fmt.Println("フォームデータの取得に失敗しました", err)
			c.JSON(500, gin.H{"error": "フォームデータの取得に失敗しました"})
			return
		}

		images := form.File["images"] // "images" はフロントエンドで設定されたフィールド名

		var urls []string
		for _, image := range images {
			url, err := cloudflare.Upload(c, image)
			if err != nil {
				fmt.Println("画像のアップロードに失敗しました", err)
				c.JSON(500, gin.H{"error": "画像のアップロードに失敗しました"})
				return
			}
			urls = append(urls, url)
		}

		c.JSON(200, gin.H{"urls": urls})
	})
}
