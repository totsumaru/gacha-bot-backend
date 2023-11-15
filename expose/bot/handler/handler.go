package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/totsumaru/gacha-bot-backend/expose/bot/handler/interaction_craete"
	"github.com/totsumaru/gacha-bot-backend/expose/bot/handler/message_create"
)

// メッセージが作成された時のハンドラです
func Handler(s *discordgo.Session) {
	s.AddHandler(message_create.MessageCreateHandler)
	s.AddHandler(interaction_craete.InteractionCreateHandler)
}
