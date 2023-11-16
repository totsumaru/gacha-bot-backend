package stripe

import (
	"encoding/json"

	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// ストライプです
type Stripe struct {
	customerID     CustomerID
	subscriptionID SubscriptionID
}

// ストライプを作成します
func NewStripe(cusID CustomerID, subID SubscriptionID) (Stripe, error) {
	s := Stripe{
		customerID:     cusID,
		subscriptionID: subID,
	}

	if err := s.validate(); err != nil {
		return s, errors.NewError("検証に失敗しました", err)
	}

	return s, nil
}

// カスタマーIDを取得します
func (s Stripe) CustomerID() CustomerID {
	return s.customerID
}

// サブスクリプションIDを取得します
func (s Stripe) SubscriptionID() SubscriptionID {
	return s.subscriptionID
}

// 検証します
func (s Stripe) validate() error {
	return nil
}

// JSONに変換します
func (s Stripe) MarshalJSON() ([]byte, error) {
	data := struct {
		CustomerID     CustomerID     `json:"customer_id"`
		SubscriptionID SubscriptionID `json:"subscription_id"`
	}{
		CustomerID:     s.customerID,
		SubscriptionID: s.subscriptionID,
	}

	return json.Marshal(data)
}

// JSONから復元します
func (s *Stripe) UnmarshalJSON(b []byte) error {
	data := struct {
		CustomerID     CustomerID     `json:"customer_id"`
		SubscriptionID SubscriptionID `json:"subscription_id"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return errors.NewError("JSONから復元できません", err)
	}

	s.customerID = data.CustomerID
	s.subscriptionID = data.SubscriptionID

	return nil
}
