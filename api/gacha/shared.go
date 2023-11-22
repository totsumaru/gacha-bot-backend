package gacha

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/application/gacha"
	domainGacha "github.com/totsumaru/gacha-bot-backend/domain/gacha"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed/button"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/result"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ========================
// ガチャのAPIの共通処理です
// ========================

// ガチャのレスポンスです
//
// リクエストにIDを追加しています。
type GachaRes struct {
	ID       string             `json:"id"`
	ServerID string             `json:"server_id"`
	Panel    EmbedReq           `json:"panel"`
	Open     EmbedReq           `json:"open"`
	Result   []ResultReq        `json:"result"`
	Role     []RoleWithPointRes `json:"role"`
	AllRole  []RoleRes          `json:"all_role"`
}

// ロールのレスポンスです
type RoleRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ロールのレスポンスです(Pointあり)
type RoleWithPointRes struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Point int    `json:"point"`
}

// ガチャのリクエストです
type GachaReq struct {
	ID       string      `json:"id"`
	ServerID string      `json:"server_id"`
	Panel    EmbedReq    `json:"panel"`
	Open     EmbedReq    `json:"open"`
	Result   []ResultReq `json:"result"`
	Role     []RoleReq   `json:"role"`
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

// ロールのリクエストです
type RoleReq struct {
	ID    string `json:"id"`
	Point int    `json:"point"`
}

// =============================================
// APIのリクエスト -> Appのリクエスト に変換します
// =============================================

// APIのリクエストをAppのリクエストに変換します
func ConvertToAppGachaReq(apiGachaReq GachaReq) gacha.GachaReq {
	panel := convertToAppEmbedReq(apiGachaReq.Panel)
	open := convertToAppEmbedReq(apiGachaReq.Open)

	var results []gacha.ResultReq
	for _, apiResult := range apiGachaReq.Result {
		results = append(results, convertToAppResultReq(apiResult))
	}

	var roles []gacha.RoleReq
	for _, apiRole := range apiGachaReq.Role {
		role := gacha.RoleReq{
			ID:    apiRole.ID,
			Point: apiRole.Point,
		}
		roles = append(roles, role)
	}

	return gacha.GachaReq{
		ID:       apiGachaReq.ID,
		ServerID: apiGachaReq.ServerID,
		Panel:    panel,
		Open:     open,
		Result:   results,
		Role:     roles,
	}
}

// 結果のリクエストをAppのリクエストに変換します
func convertToAppResultReq(apiResultReq ResultReq) gacha.ResultReq {
	appEmbedReq := convertToAppEmbedReq(apiResultReq.Embed)

	return gacha.ResultReq{
		Embed:       appEmbedReq,
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

// ====================================
// ドメイン -> APIレスポンス に変換します
// ====================================

// ガチャのドメインをAPIレスポンスに変換します
func ConvertToAPIGachaRes(
	s *discordgo.Session,
	serverID string,
	domainGacha domainGacha.Gacha,
) (GachaRes, error) {
	panel := convertDomainEmbedToAPIRes(domainGacha.Panel())
	open := convertDomainEmbedToAPIRes(domainGacha.Open())

	results := make([]ResultReq, 0)
	for _, domainResult := range domainGacha.Result() {
		results = append(results, convertDomainResultToAPIRes(domainResult))
	}

	roles := make([]RoleWithPointRes, 0)
	for _, domainRole := range domainGacha.Role() {
		id := domainRole.ID().String()

		// ロール名を取得
		r, err := s.State.Role(domainGacha.ServerID().String(), id)
		if err != nil {
			return GachaRes{}, errors.NewError("ロール名の取得に失敗しました", err)
		}

		role := RoleWithPointRes{
			ID:    id,
			Name:  r.Name,
			Point: domainRole.Point().Int(),
		}
		roles = append(roles, role)
	}

	// そのサーバーの全てのロールを取得
	guildRoles, err := s.GuildRoles(serverID)
	if err != nil {
		return GachaRes{}, errors.NewError("ロールの取得に失敗しました", err)
	}

	allRoles := make([]RoleRes, 0)
	for _, r := range guildRoles {
		role := RoleRes{
			ID:   r.ID,
			Name: r.Name,
		}
		allRoles = append(allRoles, role)
	}

	return GachaRes{
		ID:       domainGacha.ID().String(),
		ServerID: domainGacha.ServerID().String(),
		Panel:    panel,
		Open:     open,
		Result:   results,
		Role:     roles,
		AllRole:  allRoles,
	}, nil
}

// Resultのドメインをレスポンスに変換します
func convertDomainResultToAPIRes(domainResult result.Result) ResultReq {
	embedAPIRes := convertDomainEmbedToAPIRes(domainResult.Embed())

	return ResultReq{
		Embed:       embedAPIRes,
		Point:       domainResult.Point().Int(),
		Probability: domainResult.Probability().Int(),
	}
}

// Embedのドメインをレスポンスに変換します
func convertDomainEmbedToAPIRes(domainEmbed embed.Embed) EmbedReq {
	var btns []ButtonReq
	for _, v := range domainEmbed.Button() {
		btns = append(btns, convertDomainButtonToAPIRes(v))
	}

	return EmbedReq{
		Title:        domainEmbed.Title().String(),
		Description:  domainEmbed.Description().String(),
		Color:        domainEmbed.Color().Int(),
		ImageURL:     domainEmbed.ImageURL().String(),
		ThumbnailURL: domainEmbed.ThumbnailURL().String(),
		Button:       btns,
	}
}

// ボタンのドメインをレスポンスに変換します
func convertDomainButtonToAPIRes(domainButton button.Button) ButtonReq {
	return ButtonReq{
		Kind:  domainButton.Kind().String(),
		Label: domainButton.Label().String(),
		Style: domainButton.Style().String(),
	}
}
