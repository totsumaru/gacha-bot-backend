package gacha

import (
	"github.com/totsumaru/gacha-bot-backend/app/gacha"
)

// ========================
// ガチャのAPIの共通処理です
// ========================

// ガチャのリクエストです
type GachaReq struct {
	ServerID string      `json:"server_id"`
	Panel    EmbedReq    `json:"panel"`
	Open     EmbedReq    `json:"open"`
	Result   []ResultReq `json:"result"`
}

// 結果のリクエストです
type ResultReq struct {
	Embed       EmbedReq `json:"embed"`
	Point       int      `json:"point"`
	Probability int      `json:"probability"`
}

// ガチャの埋め込みのリクエストです
type EmbedReq struct {
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	Color        int         `json:"color"`
	ImageURL     string      `json:"image_url"`
	ThumbnailURL string      `json:"thumbnail_url"`
	Button       []ButtonReq `json:"button"`
}

// ボタンのリクエストです
type ButtonReq struct {
	Kind  string `json:"kind"`
	Label string `json:"label"`
	Style string `json:"style"`
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
		ServerID: apiGachaReq.ServerID,
		Panel:    panel,
		Open:     open,
		Result:   results,
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
