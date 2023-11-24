package servers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/totsumaru/gacha-bot-backend/application/server"
	domainServer "github.com/totsumaru/gacha-bot-backend/domain/server"
	"github.com/totsumaru/gacha-bot-backend/lib/auth"
	"github.com/totsumaru/gacha-bot-backend/lib/discord"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// レスポンスのサーバーです
type ServerRes struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IconURL    string `json:"icon_url"`
	Subscriber struct {
		UserName  string `json:"user_name"`
		AvatarURL string `json:"avatar_url"`
	} `json:"subscriber"`
}

// サーバーの情報を取得します
func GetAdminServers(e *gin.Engine, db *gorm.DB) {
	e.GET("/api/admin/servers", func(c *gin.Context) {
		authHeader := c.GetHeader(auth.HeaderAuthorization)

		// verify
		{
			headerRes, err := auth.GetAuthHeader(authHeader)
			if err != nil {
				errors.HandleError(c, 401, "トークンの認証に失敗しました", err)
				return
			}

			if headerRes.DiscordID != discord.TotsumaruID {
				errors.HandleError(c, 401, "アプリケーション管理者ではありません", err)
				return
			}
		}

		res := make([]ServerRes, 0)
		err := db.Transaction(func(tx *gorm.DB) error {
			sv, err := server.FindAll(tx)
			if err != nil {
				return errors.NewError("全てのサーバーを取得できません", err)
			}

			for _, s := range sv {
				r, err := convertDomainServerToApiRes(s)
				if err != nil {
					return errors.NewError("サーバーの情報を取得できません", err)
				}
				res = append(res, r)
			}

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
func convertDomainServerToApiRes(domainSv domainServer.Server) (ServerRes, error) {
	u, err := discord.Session.User(domainSv.SubscriberID().String())
	if err != nil {
		return ServerRes{}, errors.NewError("ユーザーの情報を取得できません", err)
	}

	s, err := discord.Session.Guild(domainSv.ID().String())
	if err != nil {
		return ServerRes{}, errors.NewError("サーバーの情報を取得できません", err)
	}

	res := ServerRes{}
	res.ID = domainSv.ID().String()
	res.Name = s.Name
	res.IconURL = s.IconURL("")
	res.Subscriber.UserName = u.Username
	res.Subscriber.AvatarURL = u.AvatarURL("")

	return res, nil
}
