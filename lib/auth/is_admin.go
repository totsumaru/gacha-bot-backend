package auth

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/lib/auth/permission"
	"github.com/totsumaru/gacha-bot-backend/lib/discord"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
)

// そのユーザーがサーバーの操作権限を持っていることを確認します
//
// 以下のどれか1つに当てはまる場合はtrueが返ります。
// - サーバーオーナー
// - 管理者権限ロール
func IsAdmin(serverID, userID string) error {
	if userID == discord.TotsumaruID {
		return nil
	}

	s := discord.Session

	member, err := s.GuildMember(serverID, userID)
	if err != nil {
		return errors.NewError("指定したサーバーでユーザーを取得できません", err)
	}

	guild, err := s.Guild(serverID)
	if err != nil {
		return errors.NewError("サーバーIDでギルドを取得できません", err)
	}

	// そのサーバーのオーナーの場合はtrueを返す
	if guild.OwnerID == userID {
		return nil
	}

	// サーバーにユーザーが存在している場合、操作ロールの有無を確認します
	if member.User.ID == userID {
		// サーバーの全てのロールを取得
		guildRoles, err := s.GuildRoles(serverID)
		if err != nil {
			return errors.NewError("サーバーのロール一覧を取得できません", err)
		}

		roleMap := make(map[string]*discordgo.Role)
		for _, role := range guildRoles {
			roleMap[role.ID] = role
		}

		// 管理者ロールを持っているか確認します
		for _, roleID := range member.Roles {
			// そのロールが管理者権限を持っている場合はtrueを返します
			if permission.HasPermission(roleMap[roleID].Permissions, discordgo.PermissionAdministrator) {
				return nil
			}
		}
	}

	return errors.NewError("認証に失敗しました")
}
