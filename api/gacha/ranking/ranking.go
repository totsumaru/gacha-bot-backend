package ranking

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	appUserData "github.com/totsumaru/gacha-bot-backend/application/user_data"
	domainUserData "github.com/totsumaru/gacha-bot-backend/domain/user_data"
	"github.com/totsumaru/gacha-bot-backend/lib/discord"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// レスポンスのユーザーデータです
type ResUserData struct {
	UserName  string
	AvatarURL string
	Point     int
	Rank      int
}

// ランキングを取得します
//
// 1-100位までを取得します。
func GetRanking(e *gin.Engine, db *gorm.DB) {
	e.GET("/api/gacha/ranking", func(c *gin.Context) {
		serverID := c.Query("server_id")

		res := make([]ResUserData, 0)
		err := db.Transaction(func(tx *gorm.DB) error {
			userDatas, err := appUserData.FindTop100ByServerID(tx, serverID)
			if err != nil {
				return errors.NewError("ガチャの取得に失敗しました", err)
			}

			res, err = ConvertToAPIGachaRes(discord.Session, userDatas)
			if err != nil {
				return errors.NewError("APIのレスポンスに変換できません", err)
			}

			return nil
		})
		if err != nil {
			errors.HandleError(c, 500, "ガチャを取得できません", err)
			return
		}

		c.JSON(http.StatusOK, res)
	})
}

// ドメインのユーザーデータをAPIのレスポンスに変換します
func ConvertToAPIGachaRes(
	s *discordgo.Session,
	datas []domainUserData.UserData,
) ([]ResUserData, error) {
	res := make([]ResUserData, 0)

	for i, userData := range datas {
		userID := userData.UserID().String()
		u, err := s.User(userID)
		if err != nil {
			return nil, errors.NewError("ユーザー名を取得できません", err)
		}

		resUserData := ResUserData{
			UserName:  u.Username,
			AvatarURL: u.AvatarURL(""),
			Point:     userData.Point().Int(),
			Rank:      i + 1,
		}

		res = append(res, resUserData)
	}

	return res, nil
}