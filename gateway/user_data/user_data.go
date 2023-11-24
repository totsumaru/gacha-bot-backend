package user_data

import (
	"encoding/json"
	defaultError "errors"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/user_data"
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

// Upsertは、指定されたIDに対応するレコードを更新するか、
// 存在しない場合は新しいレコードを作成します。
func (g Gateway) Upsert(userData user_data.UserData) error {
	dbUserData, err := castToDBStruct(userData)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// Upsert処理を実行
	result := g.tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // IDをキーとして使用
		UpdateAll: true,                          // 衝突した場合、すべての列を更新
	}).Create(&dbUserData)
	if result.Error != nil {
		return errors.NewError("レコードの更新/作成に失敗しました", result.Error)
	}

	return nil
}

// IDでユーザーデータを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByID(id user_data.ID) (user_data.UserData, error) {
	var res user_data.UserData

	var dbUserData gateway.UserData
	if err := g.tx.First(&dbUserData, "id = ?", id.String()).Error; err != nil {
		if err != nil {
			if defaultError.Is(err, gorm.ErrRecordNotFound) {
				// レコードが存在しない場合、NotFoundErrorを返します
				return res, errors.NotFoundError{}
			}
			return res, errors.NewError("IDでユーザーデータを取得できません", err)
		}
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbUserData)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FOR UPDATEでユーザーデータを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id user_data.ID) (user_data.UserData, error) {
	var res user_data.UserData

	var dbUserData gateway.UserData
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbUserData, "id = ?", id.String(),
	).Error; err != nil {
		if defaultError.Is(err, gorm.ErrRecordNotFound) {
			// レコードが存在しない場合、NotFoundErrorを返します
			return res, errors.NotFoundError{}
		}
		return res, errors.NewError("IDでユーザーデータを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbUserData)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// 削除します
func (g Gateway) Delete(id user_data.ID) error {
	// IDに基づいてレコードを削除
	result := g.tx.Delete(&gateway.UserData{}, "id = ?", id.String())
	if result.Error != nil {
		return errors.NewError("削除できません", result.Error)
	}

	// 主キー制約違反を検出（指定されたIDのレコードが存在しない場合）
	if result.RowsAffected == 0 {
		return errors.NewError("レコードが存在しません")
	}

	return nil
}

// FindTop100ByServerID は、指定されたサーバーIDに一致するユーザーデータの中から、
// ポイントの上位100個のレコードを取得します。
func (g Gateway) FindTop100ByServerID(serverID domain.DiscordID) ([]user_data.UserData, error) {
	var dbUserDatas []gateway.UserData
	var userDatas []user_data.UserData

	// id のプレフィックスが server_id と一致するレコードを検索し、ポイントで降順にソートして上位100件を取得
	serverIDPrefix := serverID.String() + "%" // "%" はワイルドカードとして機能します
	if err := g.tx.Where(
		"id LIKE ?", serverIDPrefix,
	).Order("point DESC").Limit(100).Find(&dbUserDatas).Error; err != nil {
		return nil, errors.NewError("データの取得に失敗しました", err)
	}

	// 取得したDBデータをドメインモデルに変換
	for _, dbUserData := range dbUserDatas {
		userData, err := castToDomainModel(dbUserData)
		if err != nil {
			return nil, errors.NewError("データの変換に失敗しました", err)
		}
		userDatas = append(userDatas, userData)
	}

	return userDatas, nil
}

// =============
// private
// =============

// ドメインモデルをDBの構造体に変換します
func castToDBStruct(userData user_data.UserData) (gateway.UserData, error) {
	dbUserData := gateway.UserData{}

	b, err := json.Marshal(&userData)
	if err != nil {
		return dbUserData, errors.NewError("Marshalに失敗しました", err)
	}

	dbUserData.ID = userData.ID().String()
	dbUserData.Data = b
	dbUserData.ServerID = userData.ServerID().String()
	dbUserData.Point = userData.Point().Int()

	return dbUserData, nil
}

// DBの構造体からドメインモデルに変換します
func castToDomainModel(dbUserData gateway.UserData) (user_data.UserData, error) {
	var res user_data.UserData

	if err := json.Unmarshal(dbUserData.Data, &res); err != nil {
		return res, errors.NewError("Unmarshalに失敗しました", err)
	}

	return res, nil
}
