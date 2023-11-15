package server

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/domain"
	"github.com/totsumaru/gacha-bot-backend/domain/server/stripe"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// サーバーです
type Server struct {
	id           domain.DiscordID
	subscriberID domain.DiscordID // 支払い者のユーザーID
	stripe       stripe.Stripe
}

// サーバーを作成します
func NewServer(
	id domain.DiscordID,
	subscriberID domain.DiscordID,
	stripe stripe.Stripe,
) (Server, error) {
	s := Server{
		id:           id,
		subscriberID: subscriberID,
		stripe:       stripe,
	}

	if err := s.validate(); err != nil {
		return Server{}, errors.NewError("サーバーの生成に失敗しました", err)
	}

	return s, nil
}

// 支払い者のIDを更新します
func (s *Server) UpdateSubscriberID(subscriberID domain.DiscordID) error {
	s.subscriberID = subscriberID
	if err := s.validate(); err != nil {
		return errors.NewError("サーバーの更新に失敗しました", err)
	}

	return nil
}

// Stripeを更新します
func (s *Server) UpdateStripe(stripe stripe.Stripe) error {
	s.stripe = stripe
	if err := s.validate(); err != nil {
		return errors.NewError("サーバーの更新に失敗しました", err)
	}

	return nil
}

// IDを返します
func (s Server) ID() domain.DiscordID {
	return s.id
}

// 支払い者のユーザーIDを返します
func (s Server) SubscriberID() domain.DiscordID {
	return s.subscriberID
}

// Stripeを返します
func (s Server) Stripe() stripe.Stripe {
	return s.stripe
}

// サーバーを検証します
func (s Server) validate() error {
	return nil
}

// サーバーをJSONに変換します
func (s Server) MarshalJSON() ([]byte, error) {
	data := struct {
		ID           domain.DiscordID `json:"id"`
		SubscriberID domain.DiscordID `json:"subscriber_id"`
		Stripe       stripe.Stripe    `json:"stripe"`
	}{
		ID:           s.id,
		SubscriberID: s.subscriberID,
		Stripe:       s.stripe,
	}

	return json.Marshal(data)
}

// JSONからサーバーに変換します
func (s *Server) UnmarshalJSON(b []byte) error {
	data := struct {
		ID           domain.DiscordID `json:"id"`
		SubscriberID domain.DiscordID `json:"subscriber_id"`
		Stripe       stripe.Stripe    `json:"stripe"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONの変換に失敗しました", err)
	}

	s.id = data.ID
	s.subscriberID = data.SubscriberID
	s.stripe = data.Stripe

	return nil
}
