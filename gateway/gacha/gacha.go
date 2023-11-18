package gacha

import (
	"encoding/json"
	defaultError "errors"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/gacha"
	"github.com/totsumaru/gacha-bot-backend/gateway"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// Upsertは、指定されたIDに対応するレコードを更新するか、
// 存在しない場合は新しいレコードを作成します。
func (g Gateway) Upsert(gacha gacha.Gacha) error {
	dbGacha, err := castToDBStruct(gacha)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// Upsert処理を実行
	result := g.tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // IDをキーとして使用
		UpdateAll: true,                          // 衝突した場合、すべての列を更新
	}).Create(&dbGacha)
	if result.Error != nil {
		return errors.NewError("レコードの更新/作成に失敗しました", result.Error)
	}

	return nil
}

// IDでガチャを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByID(id domain.UUID) (gacha.Gacha, error) {
	var res gacha.Gacha

	var dbGacha gateway.Gacha
	err := g.tx.First(&dbGacha, "id = ?", id.String()).Error
	if err != nil {
		if defaultError.Is(err, gorm.ErrRecordNotFound) {
			// レコードが存在しない場合、NotFoundErrorを返します
			return res, errors.NotFoundError{}
		}
		return res, errors.NewError("IDでガチャを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err = castToDomainModel(dbGacha)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// サーバーIDでガチャを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByServerID(serverID domain.DiscordID) (gacha.Gacha, error) {
	var res gacha.Gacha

	var dbGacha gateway.Gacha
	err := g.tx.First(&dbGacha, "server_id = ?", serverID.String()).Error
	if err != nil {
		if defaultError.Is(err, gorm.ErrRecordNotFound) {
			// レコードが存在しない場合、NotFoundErrorを返します
			return res, errors.NotFoundError{}
		}
		return res, errors.NewError("IDでガチャを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err = castToDomainModel(dbGacha)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FOR UPDATEでガチャを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id domain.DiscordID) (gacha.Gacha, error) {
	var res gacha.Gacha

	var dbGacha gateway.Gacha
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbGacha, "id = ?", id.String(),
	).Error; err != nil {
		return res, errors.NewError("IDでガチャを取得できません", err)
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
	dbGacha.ServerID = gacha.ServerID().String()
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
