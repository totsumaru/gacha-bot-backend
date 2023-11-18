package domain

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// IDです
type UUID struct {
	value string
}

// IDを作成します
func NewUUID() (UUID, error) {
	res := UUID{}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		return res, errors.NewError("UUIDの生成に失敗しました", err)
	}

	res.value = newUUID.String()

	return res, nil
}

// IDを復元します
func RestoreUUID(id string) (UUID, error) {
	res := UUID{
		value: id,
	}

	if err := res.validate(); err != nil {
		return res, errors.NewError("検証に失敗しました", err)
	}

	return res, nil
}

// IDを取得します
func (u UUID) String() string {
	return u.value
}

// IDが存在しているか確認します
func (u UUID) IsEmpty() bool {
	return u.value == ""
}

// IDを検証します
func (u UUID) validate() error {
	_, err := uuid.Parse(u.value)
	if err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}

// 構造体からJSONに変換します
func (u UUID) MarshalJSON() ([]byte, error) {
	data := struct {
		Value string `json:"value"`
	}{
		Value: u.value,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewError("Marshalに失敗しました", err)
	}

	return b, nil
}

// JSONから構造体に変換します
func (u *UUID) UnmarshalJSON(b []byte) error {
	var data struct {
		Value string `json:"value"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	u.value = data.Value

	return nil
}
