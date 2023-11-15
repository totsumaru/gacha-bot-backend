package embed

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha/embed/button"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// Embedのモデルです
type Embed struct {
	title        Title
	description  Description
	color        ColorCode
	imageURL     domain.URL
	thumbnailURL domain.URL
	button       []button.Button
}

// Embedを生成します
func NewEmbed(
	title Title,
	description Description,
	color ColorCode,
	imageURL domain.URL,
	thumbnailURL domain.URL,
	button []button.Button,
) (Embed, error) {
	e := Embed{
		title:        title,
		description:  description,
		color:        color,
		imageURL:     imageURL,
		thumbnailURL: thumbnailURL,
		button:       button,
	}

	if err := e.validate(); err != nil {
		return Embed{}, errors.NewError("Embedの生成に失敗しました", err)
	}

	return e, nil
}

// タイトルを返します
func (e Embed) Title() Title {
	return e.title
}

// 説明を返します
func (e Embed) Description() Description {
	return e.description
}

// カラーコードを返します
func (e Embed) Color() ColorCode {
	return e.color
}

// 画像URLを返します
func (e Embed) ImageURL() domain.URL {
	return e.imageURL
}

// サムネイルURLを返します
func (e Embed) ThumbnailURL() domain.URL {
	return e.thumbnailURL
}

// ボタンを返します
func (e Embed) Button() []button.Button {
	return e.button
}

// Embedを検証します
func (e Embed) validate() error {
	if len(e.button) > 5 {
		return errors.NewError("ボタンは5つまでです", nil)
	}

	return nil
}

// EmbedをJSONに変換します
func (e Embed) MarshalJSON() ([]byte, error) {
	data := struct {
		Title        Title           `json:"title"`
		Description  Description     `json:"description"`
		Color        ColorCode       `json:"color"`
		ImageURL     domain.URL      `json:"image_url"`
		ThumbnailURL domain.URL      `json:"thumbnail_url"`
		Button       []button.Button `json:"button"`
	}{
		Title:        e.title,
		Description:  e.description,
		Color:        e.color,
		ImageURL:     e.imageURL,
		ThumbnailURL: e.thumbnailURL,
		Button:       e.button,
	}

	return json.Marshal(data)
}

// JSONからEmbedを復元します
func (e *Embed) UnmarshalJSON(b []byte) error {
	data := struct {
		Title        Title           `json:"title"`
		Description  Description     `json:"description"`
		Color        ColorCode       `json:"color"`
		ImageURL     domain.URL      `json:"image_url"`
		ThumbnailURL domain.URL      `json:"thumbnail_url"`
		Button       []button.Button `json:"button"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	e.title = data.Title
	e.description = data.Description
	e.color = data.Color
	e.imageURL = data.ImageURL
	e.thumbnailURL = data.ThumbnailURL
	e.button = data.Button

	return nil
}
