package stripe

import "github.com/totsumaru/gacha-bot-backend/lib/errors"

// ストライプです
type Stripe struct {
	customerID     CustomerID
	subscriptionID SubscriptionID
}

// ストライプを作成します
func NewStripe(cusID CustomerID, subID SubscriptionID) (Stripe, error) {
	s := Stripe{}
	s.customerID = cusID
	s.subscriptionID = subID

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
