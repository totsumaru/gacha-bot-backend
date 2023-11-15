package server

import (
	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/server"
	"github.com/totsumaru/gacha-bot-backend/domain/server/stripe"
	gatewayServer "github.com/totsumaru/gacha-bot-backend/gateway/server"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/gorm"
)

// サーバーを作成します
//
// botが導入された時にコールされます。
func CreateServer(tx *gorm.DB, serverID string) error {
	id, err := domain.NewDiscordID(serverID)
	if err != nil {
		return errors.NewError("DiscordIDの生成に失敗しました", err)
	}

	emptySubscriberID := domain.DiscordID{}
	emptyStripe := stripe.Stripe{}

	s, err := server.NewServer(id, emptySubscriberID, emptyStripe)
	if err != nil {
		return errors.NewError("サーバーの生成に失敗しました", err)
	}

	gw, err := gatewayServer.NewGateway(tx)
	if err != nil {
		return errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	if err = gw.Create(s); err != nil {
		return errors.NewError("サーバーの保存に失敗しました", err)
	}

	return nil
}

// サブスクリプションを開始した時の処理です
//
// webhookのAPIからコールされます。
func StartSubscription(
	tx *gorm.DB,
	id, subscriberID, customerID, subscriptionID string,
) error {
	// idでサーバーを取得
	serverID, err := domain.NewDiscordID(id)
	if err != nil {
		return errors.NewError("IDを作成できません", err)
	}

	gw, err := gatewayServer.NewGateway(tx)
	if err != nil {
		return errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	sv, err := gw.FindByID(serverID)
	if err != nil {
		return errors.NewError("IDでサーバーを取得できません", err)
	}

	subscriver, err := domain.NewDiscordID(subscriberID)
	if err != nil {
		return errors.NewError("支払い者のIDを作成できません", err)
	}

	cusID, err := stripe.NewCustomerID(customerID)
	if err != nil {
		return errors.NewError("カスタマーIDを作成できません", err)
	}

	subsID, err := stripe.NewSubscriptionID(subscriptionID)
	if err != nil {
		return errors.NewError("サブスクリプションIDを作成できません", err)
	}

	st, err := stripe.NewStripe(cusID, subsID)
	if err != nil {
		return errors.NewError("サーバー構造体を復元できません", err)
	}

	if err = sv.UpdateSubscriberID(subscriver); err != nil {
		return errors.NewError("サーバーの更新に失敗しました", err)
	}

	if err = sv.UpdateStripe(st); err != nil {
		return errors.NewError("サーバーの更新に失敗しました", err)
	}

	if err = gw.Update(sv); err != nil {
		return errors.NewError("更新できません", err)
	}

	return nil
}

// サブスクリプションを終了した時の処理です
//
// webhookのAPIからコールされます。
// サブスクリプションの期限が終了したことを意味します。
func DeleteSubscription(tx *gorm.DB, id string) error {
	// idでサーバーを取得
	serverID, err := domain.NewDiscordID(id)
	if err != nil {
		return errors.NewError("IDを作成できません", err)
	}

	gw, err := gatewayServer.NewGateway(tx)
	if err != nil {
		return errors.NewError("ゲートウェイの生成に失敗しました", err)
	}

	s, err := gw.FindByID(serverID)
	if err != nil {
		return errors.NewError("IDでサーバーを取得できません", err)
	}

	emptySubscriberID := domain.DiscordID{}
	emptyStripe := stripe.Stripe{}

	sv, err := server.NewServer(s.ID(), emptySubscriberID, emptyStripe)
	if err != nil {
		return errors.NewError("サーバーの生成に失敗しました", err)
	}

	if err = gw.Update(sv); err != nil {
		return errors.NewError("更新できません", err)
	}

	return nil
}
