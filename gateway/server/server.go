package server

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/server"
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

// サーバーを新規作成します
//
// 同じIDのレコードが存在する場合はエラーを返します。
func (g Gateway) Create(server server.Server) error {
	dbServer, err := castToDBStruct(server)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// 新しいレコードをデータベースに保存
	result := g.tx.Create(&dbServer)
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
func (g Gateway) Update(server server.Server) error {
	dbServer, err := castToDBStruct(server)
	if err != nil {
		return errors.NewError("ドメインモデルをDBの構造体に変換できません", err)
	}

	// IDに基づいてレコードを更新
	result := g.tx.Model(&gateway.Server{}).Where(
		"id = ?",
		dbServer.ID,
	).Updates(&dbServer)
	if result.Error != nil {
		return errors.NewError("更新できません", result.Error)
	}

	// 主キー制約違反を検出（指定されたIDのレコードが存在しない場合）
	if result.RowsAffected == 0 {
		return errors.NewError("レコードが存在しません")
	}

	return nil
}

// IDでサーバーを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByID(id domain.DiscordID) (server.Server, error) {
	var res server.Server

	var dbServer gateway.Server
	if err := g.tx.First(&dbServer, "id = ?", id.String()).Error; err != nil {
		return res, errors.NewError("IDでサーバーを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbServer)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// FOR UPDATEでサーバーを取得します
//
// レコードが存在しない場合はエラーを返します。
func (g Gateway) FindByIDForUpdate(id domain.DiscordID) (server.Server, error) {
	var res server.Server

	var dbServer gateway.Server
	if err := g.tx.Set("gorm:query_option", "FOR UPDATE").First(
		&dbServer, "id = ?", id.String(),
	).Error; err != nil {
		return res, errors.NewError("IDでサーバーを取得できません", err)
	}

	// DB->ドメインモデルに変換します
	res, err := castToDomainModel(dbServer)
	if err != nil {
		return res, errors.NewError("DBをドメインモデルに変換できません", err)
	}

	return res, nil
}

// すべてのサーバーを取得します
func (g Gateway) FindAll() ([]server.Server, error) {
	var dbServers []gateway.Server
	var servers []server.Server

	// データベースからすべてのサーバーレコードを取得
	if err := g.tx.Find(&dbServers).Error; err != nil {
		return nil, errors.NewError("サーバーを取得できません", err)
	}

	// 取得したデータをドメインモデルに変換
	for _, dbServer := range dbServers {
		s, err := castToDomainModel(dbServer)
		if err != nil {
			return nil, errors.NewError("DBをドメインモデルに変換できません", err)
		}
		servers = append(servers, s)
	}

	return servers, nil
}

// 削除します
func (g Gateway) Delete(id domain.DiscordID) error {
	// IDに基づいてレコードを削除
	result := g.tx.Delete(&gateway.Server{}, "id = ?", id.String())
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
func castToDBStruct(server server.Server) (gateway.Server, error) {
	dbServer := gateway.Server{}

	b, err := json.Marshal(&server)
	if err != nil {
		return dbServer, errors.NewError("Marshalに失敗しました", err)
	}

	dbServer.ID = server.ID().String()
	dbServer.Data = b

	return dbServer, nil
}

// DBの構造体からドメインモデルに変換します
func castToDomainModel(dbServer gateway.Server) (server.Server, error) {
	var res server.Server

	if err := json.Unmarshal(dbServer.Data, &res); err != nil {
		return res, errors.NewError("Unmarshalに失敗しました", err)
	}

	return res, nil
}
