package gacha

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	"github.com/totsumaru/gacha-bot-backend/gateway"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

type Gateway struct {
	tx *gorm.DB
}

// gatewayを作成します
func NewGateway(tx *gorm.DB) (Gateway, error) {
	if tx == nil {
		return Gateway{}, errors.NewError("引数が空です")
	}

	res := Gateway{
		tx: tx,
	}

	return res, nil
}

// ガチャを新規作成します
//
// 同じIDのレコードが存在する場合はエラーを返します。
func (g Gateway) Create(gacha gacha.Gacha) error {
	dbGacha, err := castToDBStruct(gacha)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// 新しいレコードをデータベースに保存
	result := g.tx.Create(&dbGacha)
	if result.Error != nil {
		return errors.NewError("レコードを保存できませんでした", result.Error)
	}
	// 主キー制約違反を検出（同じIDのレコードが既に存在する場合）
	if result.RowsAffected == 0 {
		return errors.NewError("既存のレコードが存在しています")
	}

	return nil
}

// 更新します
func (g Gateway) Update(gacha gacha.Gacha) error {
	dbGacha, err := castToDBStruct(gacha)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// IDに基づいてレコードを更新
	result := g.tx.Model(&gateway.Gacha{}).Where(
		"id = ?",
		dbGacha.ID,
	).Updates(&dbGacha)
	if result.Error != nil {
		return errors.NewError("更新できません", result.Error)
	}

	// 主キー制約違反を検出（指定されたIDのレコードが存在しない場合）
	if result.RowsAffected == 0 {
		return errors.NewError("レコードが存在しません")
	}

	return nil
}

// IDでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByID(id domain.DiscordID) (gacha.Gacha, error) {
	var res gacha.Gacha

	var dbGacha gateway.Gacha
	if err := g.tx.First(&dbGacha, "id = ?", id.String()).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbGacha)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FOR UPDATEでアクションを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id domain.DiscordID) (gacha.Gacha, error) {
	var res gacha.Gacha

	var dbGacha gateway.Gacha
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbGacha, "id = ?", id.String(),
	).Error; err != nil {
		return res, errors.NewError("IDでアクションを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbGacha)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// 削除します
func (g Gateway) Delete(id domain.DiscordID) error {
	// IDに基づいてレコードを削除
	result := g.tx.Delete(&gateway.Gacha{}, "id = ?", id.String())
	if result.Error != nil {
		return errors.NewError("削除できません", result.Error)
	}

	// 主キー制約違反を検出（指定されたIDのレコードが存在しない場合）
	if result.RowsAffected == 0 {
		return errors.NewError("レコードが存在しません")
	}

	return nil
}

// =============
// private
// =============

// ドメインモデルをDBの構造体に変換します
func castToDBStruct(gacha gacha.Gacha) (gateway.Gacha, error) {
	dbGacha := gateway.Gacha{}

	b, err := json.Marshal(&gacha)
	if err != nil {
		return dbGacha, errors.NewError("Marshalに失敗しました", err)
	}

	dbGacha.ID = gacha.ID().String()
	dbGacha.Data = b

	return dbGacha, nil
}

// DBの構造体からドメインモデルに変換します
func castToDomainModel(dbGacha gateway.Gacha) (gacha.Gacha, error) {
	var res gacha.Gacha

	if err := json.Unmarshal(dbGacha.Data, &res); err != nil {
		return res, errors.NewError("Unmarshalに失敗しました", err)
	}

	return res, nil
}
