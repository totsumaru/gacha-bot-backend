package errors

import (
	"log"

	"github.com/gin-gonic/gin"
)

// HandleError はエラー応答を一元化して処理します。
func HandleError(c *gin.Context, statusCode int, displayMessage string, err error) {
	logContent := displayMessage
	if err != nil {
		log.Println(displayMessage + err.Error())
	}

	c.JSON(statusCode, gin.H{
		"error": logContent,
	})

	c.Abort() // これにより、その後のハンドラが実行されないようにします
}
