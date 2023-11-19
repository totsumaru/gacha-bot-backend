package count

import (
	"encoding/json"
	"time"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	libNow "github.com/totsumaru/gacha-bot-backend/lib/now"
)

// 1日の回数です
type Count struct {
	date time.Time
	num  Num
}

// 1日の上限回数を生成します
func NewCount(date time.Time, num Num) (Count, error) {
	c := Count{
		date: date,
		num:  num,
	}

	if err := c.validate(); err != nil {
		return Count{}, errors.NewError("1日の上限回数の生成に失敗しました", err)
	}

	return c, nil
}

// 日付が今日の日付であるか検証します
func (c Count) IsToday() bool {
	now := libNow.NowJST()
	return c.date.Year() == now.Year() && c.date.Month() == now.Month() && c.date.Day() == now.Day()
}

// 日付を取得します
func (c Count) Date() time.Time {
	return c.date
}

// 回数を取得します
func (c Count) Num() Num {
	return c.num
}

// 1日の上限回数を検証します
func (c Count) validate() error {
	if c.date.IsZero() {
		return errors.NewError("日付が不正です")
	}

	return nil
}

// 1日の上限回数をJSONに変換します
func (c Count) MarshalJSON() ([]byte, error) {
	data := struct {
		Date time.Time `json:"date"`
		Num  Num       `json:"num"`
	}{
		Date: c.date,
		Num:  c.num,
	}

	return json.Marshal(data)
}

// 1日の上限回数をJSONから変換します
func (c *Count) UnmarshalJSON(b []byte) error {
	data := struct {
		Date time.Time `json:"date"`
		Num  Num       `json:"num"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	c.date = data.Date
	c.num = data.Num

	return nil
}
