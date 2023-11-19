package now

import (
	"time"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(errors.NewError("時間のLocationを設定できません", err))
	}
}

// 日本時間で現在の日時を取得します
func NowJST() time.Time {
	return time.Now().In(location)
}
