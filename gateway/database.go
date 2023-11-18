package gateway

// サーバーのDBスキーマです
type Server struct {
	ID   string `gorm:"type:varchar(255);primary_key;"`
	Data []byte `gorm:"type:jsonb"`
}

// ガチャのDBスキーマです
type Gacha struct {
	ID       string `gorm:"type:uuid;primary_key;"`
	ServerID string `gorm:"type:varchar(255);not null;index:idx_server_id"`
	Data     []byte `gorm:"type:jsonb"`
}
