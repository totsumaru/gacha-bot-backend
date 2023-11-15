package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"github.com/totsumaru/gacha-bot-backend/lib/seeker"
)

type Res struct {
	DiscordID string
}

// Header(Bearer xxx)からdiscordIDを取得します
func GetAuthHeader(authHeader string) (Res, error) {
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")

	// トークンを解析
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("メソッドが期待した値と異なります: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SUPABASE_JWT_SECRET")), nil
	})
	if err != nil {
		return Res{}, errors.NewError("認証できません", err)
	}

	// トークンが有効なら、ユーザーはログインしていると判断できる
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return Res{
			DiscordID: seeker.Str(claims, []string{"user_metadata", "provider_id"}),
		}, nil
	} else {
		return Res{}, errors.NewError("トークンが無効です")
	}
}
