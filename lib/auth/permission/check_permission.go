package permission

// 指定の権限を持っているかを確認します
func HasPermission(userPermission, expectPermission int64) bool {
	return userPermission&expectPermission != 0
}
