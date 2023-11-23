package gacha

import (
	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed/button"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/result"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ガチャのリクエストです
type GachaReq struct {
	ID       string
	ServerID string
	Panel    EmbedReq
	Open     EmbedReq
	Result   []ResultReq
	Role     []RoleReq
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
	Kind     string
	Label    string
	Style    string
	URL      string
	IsHidden bool
}

// ロールのリクエストです
type RoleReq struct {
	ID    string
	Point int
}

// ========================================
// 以下、app内で使用する共通処理です
// ========================================

// 結果を作成します
func createResult(req []ResultReq) ([]result.Result, error) {
	var r []result.Result

	for _, v := range req {
		emb, err := createEmbed(v.Embed)
		if err != nil {
			return nil, errors.NewError("Embedの生成に失敗しました", err)
		}

		point, err := domain.NewPoint(v.Point)
		if err != nil {
			return nil, errors.NewError("ポイントの生成に失敗しました", err)
		}

		probability, err := result.NewProbability(v.Probability)
		if err != nil {
			return nil, errors.NewError("確率の生成に失敗しました", err)
		}

		rs, err := result.NewResult(emb, point, probability)
		if err != nil {
			return nil, errors.NewError("結果の生成に失敗しました", err)
		}

		r = append(r, rs)
	}

	return r, nil
}

// リクエストからEmbedを作成します
func createEmbed(req EmbedReq) (embed.Embed, error) {
	t, err := embed.NewTitle(req.Title)
	if err != nil {
		return embed.Embed{}, errors.NewError("タイトルの生成に失敗しました", err)
	}

	d, err := embed.NewDescription(req.Description)
	if err != nil {
		return embed.Embed{}, errors.NewError("説明の生成に失敗しました", err)
	}

	c, err := embed.NewColorCode(req.Color)
	if err != nil {
		return embed.Embed{}, errors.NewError("カラーコードの生成に失敗しました", err)
	}

	i, err := domain.NewURL(req.ImageURL)
	if err != nil {
		return embed.Embed{}, errors.NewError("画像URLの生成に失敗しました", err)
	}

	tu, err := domain.NewURL(req.ThumbnailURL)
	if err != nil {
		return embed.Embed{}, errors.NewError("サムネイルURLの生成に失敗しました", err)
	}

	b := make([]button.Button, 0)
	for _, v := range req.Button {
		kind, err := button.NewKind(v.Kind)
		if err != nil {
			return embed.Embed{}, errors.NewError("ボタンの種類の生成に失敗しました", err)
		}

		label, err := button.NewLabel(v.Label)
		if err != nil {
			return embed.Embed{}, errors.NewError("ボタンのラベルの生成に失敗しました", err)
		}

		style, err := button.NewStyle(v.Style)
		if err != nil {
			return embed.Embed{}, errors.NewError("ボタンのスタイルの生成に失敗しました", err)
		}

		url, err := domain.NewURL(v.URL)
		if err != nil {
			return embed.Embed{}, errors.NewError("URLの生成に失敗しました", err)
		}

		bb, err := button.NewButton(kind, label, style, url, v.IsHidden)
		if err != nil {
			return embed.Embed{}, errors.NewError("ボタンの生成に失敗しました", err)
		}

		b = append(b, bb)
	}

	e, err := embed.NewEmbed(t, d, c, i, tu, b)
	if err != nil {
		return embed.Embed{}, errors.NewError("埋め込みの生成に失敗しました", err)
	}

	return e, nil
}
