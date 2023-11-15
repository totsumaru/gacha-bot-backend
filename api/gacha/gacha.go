package gacha

import (
	"github.com/totsumaru/gacha-bot-backend/app/gacha"
)

// ガチャのリクエストです
type GachaReq struct {
	Panel  EmbedReq
	Open   EmbedReq
	Result []ResultReq
}

// 結果のリクエストです
type ResultReq struct {
	Embed       EmbedReq
	Point       int
	Probability int
}

// ガチャの埋め込みのリクエストです
type EmbedReq struct {
	Title        string
	Description  string
	Color        int
	ImageURL     string
	ThumbnailURL string
	Button       []ButtonReq
}

// ボタンのリクエストです
type ButtonReq struct {
	Kind  string
	Label string
	Style string
}

// APIのリクエストをAppのリクエストに変換します
func ConvertToAppGachaReq(apiGachaReq GachaReq) gacha.GachaReq {
	panel := convertToAppEmbedReq(apiGachaReq.Panel)
	open := convertToAppEmbedReq(apiGachaReq.Open)

	var results []gacha.ResultReq
	for _, apiResult := range apiGachaReq.Result {
		results = append(results, convertToAppResultReq(apiResult))
	}

	return gacha.GachaReq{
		Panel:  panel,
		Open:   open,
		Result: results,
	}
}

// 結果のリクエストをAppのリクエストに変換します
func convertToAppResultReq(apiResultReq ResultReq) gacha.ResultReq {
	embed := convertToAppEmbedReq(apiResultReq.Embed)

	return gacha.ResultReq{
		Embed:       embed,
		Point:       apiResultReq.Point,
		Probability: apiResultReq.Probability,
	}
}

// 埋め込みのリクエストをAppのリクエストに変換します
func convertToAppEmbedReq(apiEmbedReq EmbedReq) gacha.EmbedReq {
	var btns []gacha.ButtonReq
	for _, apiBtn := range apiEmbedReq.Button {
		btns = append(btns, convertToAppButtonReq(apiBtn))
	}

	return gacha.EmbedReq{
		Title:        apiEmbedReq.Title,
		Description:  apiEmbedReq.Description,
		Color:        apiEmbedReq.Color,
		ImageURL:     apiEmbedReq.ImageURL,
		ThumbnailURL: apiEmbedReq.ThumbnailURL,
		Button:       btns,
	}
}

// ボタンのリクエストをAppのリクエストに変換します
func convertToAppButtonReq(apiButtonReq ButtonReq) gacha.ButtonReq {
	return gacha.ButtonReq{
		Kind:  apiButtonReq.Kind,
		Label: apiButtonReq.Label,
		Style: apiButtonReq.Style,
	}
}
